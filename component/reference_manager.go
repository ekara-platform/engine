package component

import (
	"fmt"
	"log"

	"github.com/ekara-platform/model"
)

type (
	ReferenceManager struct {
		l  *log.Logger
		cm *ComponentManager
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

	ParentRef struct {
		Component            model.Component
		ReferencedComponents *model.ReferencedComponents
	}
)

//CreateReferenceManager creates a new references manager
func CreateReferenceManager(cm *ComponentManager) *ReferenceManager {
	return &ReferenceManager{
		l:                       cm.Logger,
		cm:                      cm,
		UsedReferences:          model.CreateUsedReferences(),
		Orphans:                 model.CreateOrphans(),
		ReferencedComponents:    model.CreateReferencedComponents(),
		Parents:                 make([]ParentRef, 0, 0),
		SortedFetchedComponents: make([]string, 0, 0),
		rootComponents:          model.CreateReferencedComponents(),
	}
}

func (rm *ReferenceManager) Init(c model.Component) error {
	rm.l.Println("Parsing the main descriptor")
	rm.rootComponent = c
	fComp, err := fetch(rm.cm, c)
	if err != nil {
		return err
	}
	rm.cm.Paths[c.Id] = fComp
	refs, err := model.ParseYamlDescriptorReferences(fComp.DescriptorUrl, rm.cm.data)
	if err != nil {
		return err
	}

	// Keep all the references used into the main descriptor
	u, o := refs.Uses(rm.Orphans)
	rm.UsedReferences.AddAll(*u)
	rm.Orphans.AddAll(*o)
	referenced, err := refs.References(c.Id)
	if err != nil {
		return err
	}
	for _, v := range referenced.Refs {
		// Keep all the component declared into the main descriptor
		rm.ReferencedComponents.AddReference(v)
		rm.l.Printf("Adding %s as component into the main descriptor", v.Component.Id)
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
	err = rm.parseParent(parentC)

	return err
}

func (rm *ReferenceManager) parseParent(p model.Component) error {
	rm.l.Printf("Parsing parent %s", p.Repository)

	p.Id = fmt.Sprintf("%s%d", p.Id, len(rm.Parents)+1)
	rm.l.Printf("Parent Renamed %s", p.Id)

	fComp, err := fetch(rm.cm, p)
	if err != nil {
		rm.l.Printf("error fetching the parent %s", err.Error())
		return err
	}

	refs, err := model.ParseYamlDescriptorReferences(fComp.DescriptorUrl, rm.cm.data)
	if err != nil {
		rm.l.Printf("error parsing the parent references %s", err.Error())
		return err
	}
	// Keep all the references used into the parent
	u, o := refs.Uses(rm.Orphans)
	rm.UsedReferences.AddAll(*u)
	rm.Orphans.AddAll(*o)

	pRef := ParentRef{
		Component:            p,
		ReferencedComponents: model.CreateReferencedComponents(),
	}
	// Keep only the references of the new declaration of compoments
	referenced, err := refs.References(p.Id)
	if err != nil {
		return err
	}
	for _, v := range referenced.Refs {
		rm.l.Printf("Component referenced by the parent %s", v.Component.Repository)
		added := rm.ReferencedComponents.AddReference(v)
		if added {
			rm.l.Println("Referenced added")
			pRef.ReferencedComponents.AddReference(v)
		} else {
			rm.l.Println("Referenced ignored")
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
		return rm.parseParent(parentC)
	}
	return nil
}

func (rm *ReferenceManager) Ensure() error {
	rm.l.Printf("Reference manager ensuring, main descriptor %s", rm.rootComponent.Id)
	// Root component is the main descriptor
	err := rm.callEnsure(rm.rootComponent)
	if err != nil {
		return err
	}
	// All the parents must be processed
	for _, p := range rm.Parents {
		rm.l.Printf("Reference manager, ensuring parent %s", p.Component.Id)
		err := rm.callEnsure(p.Component)
		if err != nil {
			return err
		}

		// Once a parent has been processed	we must process its declared
		// components in alphabetical order
		for _, c := range p.ReferencedComponents.Sorted() {
			rm.l.Printf("Reference manager, ensuring parent component %s", c.Id)
			if rm.UsedReferences.IdUsed(c.Id) {
				err := rm.callEnsure(c)
				if err != nil {
					return err
				}
			} else {
				rm.l.Printf("Reference manager: unused parent component %s", c.Id)
			}
		}
	}
	// Once all the parents have been processed we must process
	// all the components declared into the main descriptor, this
	// is alo done in alphabetical order
	for _, c := range rm.rootComponents.Sorted() {
		rm.l.Printf("Reference manager, ensuring descriptor component %s", c.Id)
		if rm.UsedReferences.IdUsed(c.Id) {
			err := rm.callEnsure(c)
			if err != nil {
				return err
			}
		} else {
			rm.l.Printf("Reference manager: unused descriptor component %s", c.Id)
		}
	}
	return nil
}

func (rm *ReferenceManager) callEnsure(c model.Component) error {
	rm.cm.Environment().Platform().AddComponent(c)
	err := rm.cm.EnsureOneComponent(c)
	if err != nil {
		return err
	}
	rm.SortedFetchedComponents = append(rm.SortedFetchedComponents, c.Id)
	return nil
}
