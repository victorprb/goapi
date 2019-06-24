package requests

type errors map[string][]string

func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}
