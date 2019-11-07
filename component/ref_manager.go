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
		usedReferences *model.UsedReferences
		orphans        *model.Orphans
		//referencedComponents stores the references of all components declared
		// into the parsed descriptors and all its parents
		referencedComponents    *model.ReferencedComponents
		parents                 []ParentRef
		rootComponent           model.Component
		rootComponents          *model.ReferencedComponents
		sortedFetchedComponents []string
	}
	//ParentRef defines the reference to the descriptor's parent
	ParentRef struct {
		comp                 model.Component
		referencedComponents *model.ReferencedComponents
	}
)

//createReferenceManager creates a new references manager
func createReferenceManager(l *log.Logger) *ReferenceManager {
	return &ReferenceManager{
		l:                       l,
		usedReferences:          model.CreateUsedReferences(),
		orphans:                 model.CreateOrphans(),
		referencedComponents:    model.CreateReferencedComponents(),
		parents:                 make([]ParentRef, 0, 0),
		sortedFetchedComponents: make([]string, 0, 0),
		rootComponents:          model.CreateReferencedComponents(),
	}
}

func (rm *ReferenceManager) init(c model.Component, cm *manager, ctx *model.TemplateContext) error {
	rm.l.Println("Parsing the main descriptor")
	rm.rootComponent = c

	url, _, err := cm.ensureOneComponent(c)
	if err != nil {
		rm.l.Printf("error fetching the main descriptor %s", err.Error())
		return err
	}

	refs, err := model.ParseYamlDescriptorReferences(url, ctx)
	if err != nil {
		return err
	}

	// Keep all the references used into the main descriptor
	u, o := refs.Uses(rm.orphans)
	for id := range u.Refs {
		rm.l.Printf("Component used by the main descriptor %s", id)
		rm.usedReferences.AddReference(id)
	}
	for id := range o.Refs {
		rm.l.Printf("Orphan used by the main descriptor %s", id)
		rm.orphans.AddReference(id)
	}

	referenced, err := refs.References(c.Id)
	if err != nil {
		return err
	}
	for _, v := range referenced.Refs {
		// Keep all the component declared into the main descriptor
		rm.referencedComponents.AddReference(v)
		rm.l.Printf("Adding %s as component into the main descriptor", v.Component.Id)
		rm.rootComponents.AddReference(v)
	}

	// if exists we parse the parent of the descriptor
	parent, b, err := refs.Parent()
	if err != nil {
		return err
	}

	if b {
		parentC, _ := parent.Component()
		err = rm.parseParent(parentC, cm, ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (rm *ReferenceManager) parseParent(p model.Component, cm *manager, ctx *model.TemplateContext) error {
	rm.l.Printf("Parsing parent %s", p.Repository)

	p.Id = fmt.Sprintf("%s%d", p.Id, len(rm.parents)+1)
	rm.l.Printf("Parent Renamed %s", p.Id)

	url, _, err := cm.ensureOneComponent(p)
	if err != nil {
		rm.l.Printf("error fetching the parent %s", err.Error())
		return err
	}

	refs, err := model.ParseYamlDescriptorReferences(url, ctx)
	if err != nil {
		rm.l.Printf("error parsing the parent references %s", err.Error())
		return err
	}
	// Keep all the references used into the parent
	u, o := refs.Uses(rm.orphans)
	for id := range u.Refs {
		rm.l.Printf("Component used by the parent %s", id)
		rm.usedReferences.AddReference(id)
	}

	for id := range o.Refs {
		rm.l.Printf("Orphan used by the parent %s", id)
		rm.orphans.AddReference(id)
	}

	pRef := ParentRef{
		comp:                 p,
		referencedComponents: model.CreateReferencedComponents(),
	}
	// Keep only the references of the new declaration of compoments
	referenced, err := refs.References(p.Id)
	if err != nil {
		return err
	}
	for _, v := range referenced.Refs {
		rm.l.Printf("Component referenced by the parent %s", v.Component.Repository)
		added := rm.referencedComponents.AddReference(v)
		if added {
			rm.l.Println("Referenced added")
			pRef.referencedComponents.AddReference(v)
		} else {
			rm.l.Println("Referenced ignored")
		}
	}
	rm.parents = append(rm.parents, pRef)

	// If the descriptor has a parent then we start fetching all the parents
	parent, b, err := refs.Parent()
	if err != nil {
		return err
	}

	if b {
		parentC, _ := parent.Component()
		return rm.parseParent(parentC, cm, ctx)
	}
	return nil
}

func (rm *ReferenceManager) ensure(env *model.Environment, cm *manager, ctx *model.TemplateContext) error {
	// The parents content must be processed fist
	for i := len(rm.parents) - 1; i >= 0; i-- {
		p := rm.parents[i]

		// we must process the parent's declared components before the parent itself
		// this will be done based on the alphabetical order of the component names
		for _, c := range p.referencedComponents.Sorted() {
			rm.l.Printf("Reference manager, ensuring parent component %s", c.Id)
			if rm.usedReferences.IdUsed(c.Id) {
				err := rm.callEnsure(c, env, cm, ctx)
				if err != nil {
					return err
				}
			} else {
				rm.l.Printf("Reference manager: unused parent component %s", c.Id)
			}
		}

		//Once the declared components have been processed we can process the parent
		rm.l.Printf("Reference manager, ensuring parent %s", p.comp.Id)
		err := rm.callEnsure(p.comp, env, cm, ctx)
		if err != nil {
			return err
		}
	}

	// Once all the parents have been processed we must process
	// all the components declared into the main descriptor, this
	// is alo done in alphabetical order
	for _, c := range rm.rootComponents.Sorted() {
		rm.l.Printf("Reference manager, ensuring descriptor component %s", c.Id)
		if rm.usedReferences.IdUsed(c.Id) {
			err := rm.callEnsure(c, env, cm, ctx)
			if err != nil {
				return err
			}
		} else {
			rm.l.Printf("Reference manager: unused descriptor component %s", c.Id)
		}
	}

	rm.l.Printf("Reference manager ensuring, main descriptor %s", rm.rootComponent.Id)
	// Root component is the main descriptor
	err := rm.callEnsure(rm.rootComponent, env, cm, ctx)
	if err != nil {
		return err
	}

	return nil
}

func (rm *ReferenceManager) callEnsure(c model.Component, env *model.Environment, cm *manager, ctx *model.TemplateContext) error {
	env.Platform().AddComponent(c)

	url, hasDesc, err := cm.ensureOneComponent(c)
	if err != nil {
		return err
	}

	if hasDesc {
		rm.l.Printf("creating partial environment based on component %s", c.Id)
		descriptorYaml, err := model.ParseYamlDescriptor(url, ctx)
		if err != nil {
			rm.l.Printf("error parsing the descriptor: %s", err.Error())
			return err
		}

		cEnv, err := model.CreateEnvironment(url.String(), descriptorYaml, c.Id)
		if err != nil {
			return err
		}

		rm.l.Println("prepare partial environment for customization")
		err = env.Customize(c, cEnv)

		if err != nil {
			rm.l.Printf("error customizing the environment %s", err.Error())
			return err
		}

		ctx.Model = model.CreateTEnvironmentForEnvironment(*env)
	}
	rm.sortedFetchedComponents = append(rm.sortedFetchedComponents, c.Id)
	return nil
}

//Component returns the referenced component
func (p ParentRef) Component() (model.Component, error) {
	return p.comp, nil
}

//ComponentName returns the referenced component name
func (p ParentRef) ComponentName() string {
	return p.comp.Id
}
