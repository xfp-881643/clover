package encoding

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type BaseModel struct {
	ID string `clover:""`
}

type TestStruct struct {
	BaseModel
	IntField    int                    `clover:"int,omitempty"`
	UintField   uint                   `clover:"uint,omitempty"`
	StringField string                 `clover:",omitempty"`
	FloatField  float32                `clover:",omitempty"`
	BoolField   bool                   `clover:",omitempty"`
	TimeField   time.Time              `clover:",omitempty"`
	IntPtr      *int                   `clover:",omitempty"`
	SliceField  []int                  `clover:",omitempty"`
	MapField    map[string]interface{} `clover:",omitempty"`
	Data        []byte                 `clover:",omitempty"`
}

func TestNormalize(t *testing.T) {
	date := time.Date(2020, 0o1, 1, 0, 0, 0, 0, time.UTC)

	var x int = 100

	s := &TestStruct{
		BaseModel: BaseModel{
			ID: "UID",
		},
		TimeField:   date,
		IntField:    10,
		FloatField:  0.1,
		StringField: "aString",
		BoolField:   true,
		IntPtr:      &x,
		SliceField:  []int{1, 2, 3, 4},
		Data:        []byte("hello, clover!"),
		MapField: map[string]interface{}{
			"hello": "clover",
		},
	}

	ns, err := Normalize(s)
	require.NoError(t, err)

	require.IsType(t, ns, map[string]interface{}{})

	m := ns.(map[string]interface{})
	require.IsType(t, m["Data"], []byte{})

	require.Nil(t, m["uint"]) // testing omitempty
	require.Equal(t, m["IntPtr"], int64(100))

	require.Equal(t, m["ID"], "UID")

	s1 := &TestStruct{}
	err = Convert(m, s1)
	require.NoError(t, err)

	require.Equal(t, s1.ID, "UID")

	require.Equal(t, s, s1)

	err = Convert(m, 10)
	require.Error(t, err)
}

func TestNormalize2(t *testing.T) {
	norm, err := Normalize(nil)
	require.NoError(t, err)
	require.Nil(t, norm)

	_, err = Normalize(make(chan struct{}))
	require.Error(t, err)

	_, err = Normalize(make(map[int]interface{}))
	require.Error(t, err)
}

func TestNormalize3(t *testing.T) {
	date := time.Date(2020, 01, 1, 0, 0, 0, 0, time.UTC)

	s := &TestStruct{
		TimeField:   date,
		UintField:   0,
		IntField:    0,
		FloatField:  0,
		StringField: "",
		BoolField:   false,
		IntPtr:      nil,
		SliceField:  []int{},
		Data:        nil,
		MapField:    map[string]interface{}{},
	}

	ns, err := Normalize(s)
	require.NoError(t, err)

	require.IsType(t, ns, map[string]interface{}{})
	m := ns.(map[string]interface{})

	require.Nil(t, m["int"])
	require.Nil(t, m["uint"])
	require.Nil(t, m["FloatField"])
	require.Nil(t, m["BoolField"])
	require.Nil(t, m["SliceField"])
	require.Nil(t, m["Data"])
	require.Nil(t, m["MapField"])
	require.Nil(t, m["IntPtr"])
}

func TestEncodeDecode(t *testing.T) {
	date := time.Date(2020, 01, 1, 0, 0, 0, 0, time.UTC)

	s := &TestStruct{
		TimeField:  date,
		SliceField: []int{1, 2, 3, 4},
		MapField: map[string]interface{}{
			"hello": "clover",
		},
	}

	data, err := Encode(s)
	require.NoError(t, err)

	s1 := &TestStruct{}
	require.NoError(t, Decode(data, s1))

	require.Equal(t, s, s1)
}
