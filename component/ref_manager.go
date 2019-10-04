package component

import (
	"fmt"
	"log"

	"github.com/ekara-platform/model"
)

type (
	//ReferenceManager manages all the components declared and used into a descriptor
	ReferenceManager struct {
		l  *log.Logger
		cm *manager
		//usedReferences stores the references of all components used
		// into the parsed descriptors and all its parents
		UsedReferences *model.UsedReferences
		Orphans        *model.Orphans
		//referencedComponents stores the references of all components declared
		// into the parsed descriptors and all its parents
		ReferencedComponents    *model.ReferencedComponents
		Parents                 []ParentRef
		rootComponent           model.Component
		rootComponents          *model.ReferencedComponents
		SortedFetchedComponents []string
	}
	//ParentRef defines the reference to the descriptor's parent
	ParentRef struct {
		Comp                 model.Component
		ReferencedComponents *model.ReferencedComponents
	}
)

//CreateReferenceManager creates a new references manager
func CreateReferenceManager(cm *manager) *ReferenceManager {
	return &ReferenceManager{
		cm:                      cm,
		UsedReferences:          model.CreateUsedReferences(),
		Orphans:                 model.CreateOrphans(),
		ReferencedComponents:    model.CreateReferencedComponents(),
		Parents:                 make([]ParentRef, 0, 0),
		SortedFetchedComponents: make([]string, 0, 0),
		rootComponents:          model.CreateReferencedComponents(),
	}
}

func (rm *ReferenceManager) init(c model.Component) error {
	rm.cm.lC.Log().Println("Parsing the main descriptor")
	rm.rootComponent = c
	fComp, err := fetch(rm.cm, c)
	if err != nil {
		return err
	}
	rm.cm.Paths[c.Id] = fComp
	refs, err := model.ParseYamlDescriptorReferences(fComp.DescriptorUrl, rm.cm.TemplateContext())
	if err != nil {
		return err
	}

	// Keep all the references used into the main descriptor
	u, o := refs.Uses(rm.Orphans)
	for id := range u.Refs {
		rm.cm.lC.Log().Printf("Component used by the main descriptor %s", id)
		rm.UsedReferences.AddReference(id)
	}
	for id := range o.Refs {
		rm.cm.lC.Log().Printf("Orphan used by the main descriptor %s", id)
		rm.Orphans.AddReference(id)
	}

	referenced, err := refs.References(c.Id)
	if err != nil {
		return err
	}
	for _, v := range referenced.Refs {
		// Keep all the component declared into the main descriptor
		rm.ReferencedComponents.AddReference(v)
		rm.cm.lC.Log().Printf("Adding %s as component into the main descriptor", v.Component.Id)
		rm.rootComponents.AddReference(v)
	}

	// We always parse the parent of the descriptor
	parent, err := refs.Parent()
	if err != nil {
		return err
	}
	parentC, err := parent.Component()
	if err != nil {
		return err
	}
	err = rm.parseParent(parentC, rm.cm.TemplateContext())

	return err
}

func (rm *ReferenceManager) parseParent(p model.Component, data *model.TemplateContext) error {
	rm.cm.lC.Log().Printf("Parsing parent %s", p.Repository)

	p.Id = fmt.Sprintf("%s%d", p.Id, len(rm.Parents)+1)
	rm.cm.lC.Log().Printf("Parent Renamed %s", p.Id)

	fComp, err := fetch(rm.cm, p)
	if err != nil {
		rm.cm.lC.Log().Printf("error fetching the parent %s", err.Error())
		return err
	}

	refs, err := model.ParseYamlDescriptorReferences(fComp.DescriptorUrl, data)
	if err != nil {
		rm.cm.lC.Log().Printf("error parsing the parent references %s", err.Error())
		return err
	}
	// Keep all the references used into the parent
	u, o := refs.Uses(rm.Orphans)
	for id := range u.Refs {
		rm.cm.lC.Log().Printf("Component used by the parent %s", id)
		rm.UsedReferences.AddReference(id)
	}
	for id := range o.Refs {
		rm.cm.lC.Log().Printf("Orphan used by the parent %s", id)
		rm.Orphans.AddReference(id)
	}

	pRef := ParentRef{
		Comp:                 p,
		ReferencedComponents: model.CreateReferencedComponents(),
	}
	// Keep only the references of the new declaration of compoments
	referenced, err := refs.References(p.Id)
	if err != nil {
		return err
	}
	for _, v := range referenced.Refs {
		rm.cm.lC.Log().Printf("Component referenced by the parent %s", v.Component.Repository)
		added := rm.ReferencedComponents.AddReference(v)
		if added {
			rm.cm.lC.Log().Println("Referenced added")
			pRef.ReferencedComponents.AddReference(v)
		} else {
			rm.cm.lC.Log().Println("Referenced ignored")
		}
	}
	rm.Parents = append(rm.Parents, pRef)

	// If the descriptor has a parent then we start fetching all the parents
	if refs.Ekara.Parent.Repository != "" {
		parent, err := refs.Parent()
		if err != nil {
			return err
		}
		parentC, err := parent.Component()
		if err != nil {
			return err
		}
		return rm.parseParent(parentC, data)
	}
	return nil
}

func (rm *ReferenceManager) ensure() error {
	// The parents content must be processed fist
	for i := len(rm.Parents) - 1; i >= 0; i-- {
		p := rm.Parents[i]

		// we must process the parent's declared components before the parent itself
		// this will be done based on the alphabetical order of the component names
		for _, c := range p.ReferencedComponents.Sorted() {
			rm.cm.lC.Log().Printf("Reference manager, ensuring parent component %s", c.Id)
			if rm.UsedReferences.IdUsed(c.Id) {
				err := rm.callEnsure(c, rm.cm.TemplateContext())
				if err != nil {
					return err
				}
			} else {
				rm.cm.lC.Log().Printf("Reference manager: unused parent component %s", c.Id)
			}
		}

		//Once the declared components have been processed we can process the parent
		rm.cm.lC.Log().Printf("Reference manager, ensuring parent %s", p.Comp.Id)
		err := rm.callEnsure(p.Comp, rm.cm.TemplateContext())
		if err != nil {
			return err
		}

	}

	// Once all the parents have been processed we must process
	// all the components declared into the main descriptor, this
	// is alo done in alphabetical order
	for _, c := range rm.rootComponents.Sorted() {
		rm.cm.lC.Log().Printf("Reference manager, ensuring descriptor component %s", c.Id)
		if rm.UsedReferences.IdUsed(c.Id) {
			err := rm.callEnsure(c, rm.cm.TemplateContext())
			if err != nil {
				return err
			}
		} else {
			rm.cm.lC.Log().Printf("Reference manager: unused descriptor component %s", c.Id)
		}
	}

	rm.cm.lC.Log().Printf("Reference manager ensuring, main descriptor %s", rm.rootComponent.Id)
	// Root component is the main descriptor
	err := rm.callEnsure(rm.rootComponent, rm.cm.TemplateContext())
	if err != nil {
		return err
	}

	return nil
}

func (rm *ReferenceManager) callEnsure(c model.Component, data *model.TemplateContext) error {
	rm.cm.Environment().Platform().AddComponent(c)
	err := rm.cm.ensureOneComponent(c, data)
	if err != nil {
		return err
	}
	rm.SortedFetchedComponents = append(rm.SortedFetchedComponents, c.Id)
	return nil
}

//Component returns the referenced component
func (p ParentRef) Component() (model.Component, error) {
	return p.Comp, nil
}

//ComponentName returns the referenced component name
func (p ParentRef) ComponentName() string {
	return p.Comp.Id
}
