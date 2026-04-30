package getparamid

import "strconv"

func HandleParamIdUrl(name string) (id int, err error) {
	id, err = strconv.Atoi(name)

	if err != nil {
		return 0, err
	}

	return id, nil
}
