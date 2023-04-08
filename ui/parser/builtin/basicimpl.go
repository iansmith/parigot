//go:build js

package builtin

import (
	"log"
	"syscall/js"

	"github.com/iansmith/parigot/apiimpl/dom"
	dommsg "github.com/iansmith/parigot/g/msg/dom/v1"
	lib "github.com/iansmith/parigot/lib/go"
)

// ToggleSingle either adds or removes the arg[0] class to the
// element implied by "this".  God knows if that `this` always works in JS.
func ToggleSingle(_ js.Value, arg []js.Value) any {
	_ = arg[0] // you only really need this if you have
	// multiple elements that share this handler... and if so you can probably
	// more easily solve your problems with a css class.

	toggleClass := arg[1].String()
	selector := arg[2].String()
	serverId := arg[3].Float()
	server := dom.FindByServerId(serverId)

	// again, we don't this because we know the only object that is going
	// to receive this message is given by selector
	//	eventTarget := evt.Get("target")

	resp, err := server.ElementById(&dommsg.ElementByIdRequest{
		Id: selector[1:],
	})
	if err != nil {
		log.Printf("WARNING: element by id failed (%s): %s", selector, err.Error())
		return js.Null()
	}
	var found bool
	part := resp.Elem.GetTag().GetCssClass()
	found, part = removeElement(part, toggleClass)
	if !found {
		part = append(part, toggleClass)
	}
	resp.Elem.GetTag().CssClass = part
	err = server.UpdateCssClass(&dommsg.UpdateCssClassRequest{Elem: resp.Elem})
	if err != nil {
		log.Printf("WARN: Unable to set the CSS Classes on %s", resp.Elem.Tag.Id)
	}
	return js.Null()
}

func removeElement(part []string, candidate string) (bool, []string) {
	found := -1

	for i := 0; i < len(part); i++ {
		if part[i] == candidate {
			found = i
			break
		}
	}
	if found == -1 {
		return false, part
	}
	if found == 0 {
		if len(part) == 1 {
			part = []string{}
		} else {
			part = part[1:]
		}
	} else if found == len(part)-1 {
		if len(part)-1 == 0 {
			part = []string{}
		} else {
			part = part[:len(part)-1]
		}
	} else {
		// not the zeroth or last element
		part = append(part[:found], part[found+1:]...)
	}
	return true, part
}

// ParigotId returns the string value of the id given.
func ParigotId[T lib.AllIdPtr](id T) string {
	return lib.Unmarshal(id).String()
}
