package sortedmap

import (
	"encoding/json"
	"fmt"
)

func ExampleSortedMap_MarshalJSON() {
	type Item struct {
		Key   string
		Value string
	}

	m1 := New[string, *Item](0)
	m1.Set("y", &Item{"y", "1"})
	m1.Set("x", &Item{"x", "2"})

	data, err := json.MarshalIndent(m1, "", "  ")
	fmt.Println(string(data), err)

	m2 := New[int, *Item](0)
	m2.Set(2, &Item{"y", "2"})
	m2.Set(1, &Item{"x", "1"})

	data, err = json.MarshalIndent(m2, "", "  ")
	fmt.Println(string(data), err)

	// Output: {
	//   "y": {
	//     "Key": "y",
	//     "Value": "1"
	//   },
	//   "x": {
	//     "Key": "x",
	//     "Value": "2"
	//   }
	// } <nil>
	// {
	//   "2": {
	//     "Key": "y",
	//     "Value": "2"
	//   },
	//   "1": {
	//     "Key": "x",
	//     "Value": "1"
	//   }
	// } <nil>
}

func ExampleSortedMap_UnmarshalJSON() {
	type Item struct {
		Value int `json:"v"`
	}

	data := []byte(`{"a": {"v":1},"c": {"v":2},"b": {"v":3}}`)
	m1 := New[string, Item](0)
	_ = json.Unmarshal(data, &m1)

	data, _ = json.MarshalIndent(m1, "", "  ")
	fmt.Println(string(data))

	// ----

	m2 := New[string, *Item](0)
	_ = json.Unmarshal(data, &m2)

	data, _ = json.MarshalIndent(m2, "", "  ")
	fmt.Println(string(data))

	// ----

	data = []byte(`{"a": 1,"c": 2,"b": 3}`)
	m3 := New[string, int](0)
	_ = json.Unmarshal(data, &m3)

	data, _ = json.MarshalIndent(m3, "", "  ")
	fmt.Println(string(data))

	// Output: {
	//   "a": {
	//     "v": 1
	//   },
	//   "c": {
	//     "v": 2
	//   },
	//   "b": {
	//     "v": 3
	//   }
	// }
	// {
	//   "a": {
	//     "v": 1
	//   },
	//   "c": {
	//     "v": 2
	//   },
	//   "b": {
	//     "v": 3
	//   }
	// }
	// {
	//   "a": 1,
	//   "c": 2,
	//   "b": 3
	// }
}
