package jsongumtree

import (
	"encoding/json"
	"fmt"
	"testing"
)

func PrettyPrint(v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(string(b))
}

func TestMap(t *testing.T) {
	jsonData := []byte(`{
		"user": {
			"id": 123,
			"name": "Juan",
			"email": "juan@example.com",
			"active": true,
			"verified": false,
			"balance": 1500.50,
			"metadata": null
		},
		"hobbies": ["leer", "programar", "correr"],
		"addresses": [
			{
				"type": "home",
				"street": "Calle 123",
				"number": 456,
				"coordinates": {
					"lat": 19.4326,
					"lng": -99.1332
				}
			},
			{
				"type": "work",
				"street": "Avenida Principal",
				"number": 789,
				"coordinates": {
					"lat": 19.4284,
					"lng": -99.1276
				}
			}
		],
		"settings": {
			"notifications": {
				"email": true,
				"sms": false,
				"push": true
			},
			"privacy": {
				"profile_public": false,
				"show_email": false
			}
		},
		"empty_array": [],
		"empty_object": {},
		"mixed_array": [1, "text", true, null, {"nested": "value"}]
	}`)

	root := NewRootNode()
	if err := root.Parse(jsonData); err != nil {
		t.Fatal(err)
	}

	// Pretty print para ver el resultado
	PrettyPrint(root)
}
