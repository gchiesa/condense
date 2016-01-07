package template
import (
	"fmt"
)

type Rule func(path []interface{}, node interface{}) (newKey interface{}, newNode interface{})
type Rules struct {
	Early []Rule
	Depth []Rule
}

func (r *Rules) AttachEarly(rule Rule) {
	r.Early = append(r.Early, rule)
}

func (r *Rules) Attach(rule Rule) {
	r.Depth = append(r.Depth, rule)
}

func Walk(path []interface{}, node interface{}, rules *Rules) (newKey interface{}, newNode interface{}) {
	newPath := make([]interface{}, len(path))
	copy(newPath, path)

	newNode = node
	newKey = interface{}(nil)
	if len(newPath) > 0 {
		newKey = newPath[len(newPath)-1]
	}

	for _, rule := range rules.Early {
		newKey, newNode = rule(newPath, newNode)
		if len(newPath) > 0 {
			if newKey == nil {
				return nil, nil
			}

			newPath[len(newPath)-1] = newKey
		}
	}

	switch typed := newNode.(type) {
	default:
		panic(fmt.Sprintf("unknown type: %T\n", typed))
	case []interface{}:
		filtered := []interface{}{}
		for deepIndex, deepNode := range typed {
			newDeepPath := make([]interface{}, len(newPath)+1)
			copy(newDeepPath, newPath)
			newDeepPath[cap(newDeepPath)-1] = deepIndex

			newDeepIndex := interface{}(deepIndex)
			newDeepNode := deepNode
			newDeepIndex, newDeepNode = Walk(newDeepPath, newDeepNode, rules)

			if newDeepIndex != nil {
				filtered = append(filtered, newDeepNode)
			}
		}

		newNode = interface{}(filtered)
	case map[string]interface{}:
		filtered := make(map[string]interface{})
		for deepKey, deepNode := range typed {
			newDeepPath := make([]interface{}, len(newPath)+1)
			copy(newDeepPath, newPath)
			newDeepPath[cap(newDeepPath)-1] = deepKey

			newDeepKey := interface{}(deepKey)
			newDeepNode := deepNode
			newDeepKey, newDeepNode = Walk(newDeepPath, newDeepNode, rules)

			if newDeepKey != nil {
				filtered[newDeepKey.(string)] = newDeepNode
			}
		}

		newNode = interface{}(filtered)
	case string:
	case bool:
	case int:
	case float64:
	}

	for _, rule := range rules.Depth {
		newKey, newNode = rule(newPath, newNode)
		if len(newPath) > 0 {
			if newKey == nil {
				break
			}

			newPath[len(newPath)-1] = newKey
		}
	}

	return newKey, newNode
}

func Process(node interface{}, rules *Rules) interface{} {
	emptyPath := []interface{}{}
	_, processed := Walk(emptyPath, node, rules)
	return processed
}