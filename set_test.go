package matroid

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
)

const type1 ElementType = "TYPE1"

type testElement1 struct {
	V int
	W float64
}

func (e testElement1) GetType() ElementType {
	return type1
}

func (e testElement1) Key() string {
	return strconv.Itoa(e.V)
}

func (e testElement1) Value() interface{} {
	return e.V
}

func (e testElement1) Weight() float64 {
	return e.W
}

const type2 ElementType = "TYPE2"

type testElement2 struct {
	V int
	W float64
}

func (e testElement2) GetType() ElementType {
	return type2
}

func (e testElement2) Key() string {
	return strconv.Itoa(e.V)
}

func (e testElement2) Value() interface{} {
	return e.V
}

func (e testElement2) Weight() float64 {
	return e.W
}

func TestNewSet(t *testing.T) {
	type args struct {
		t ElementType
		e []Element
	}
	tests := []struct {
		name    string
		args    args
		want    *Set
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				t: type1,
				e: []Element{testElement1{V: 1}, testElement1{V: 2}, testElement1{V: 3}},
			},
			want: &Set{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
					"3": testElement1{V: 3},
				},
				setType: type1,
			},
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				t: type1,
				e: []Element{testElement1{V: 1}, testElement1{V: 2}, testElement2{V: 3}},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSet(tt.args.t, tt.args.e...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Add(t *testing.T) {
	type fields struct {
		set           map[string]Element
		groundSetType ElementType
	}
	type args struct {
		e Element
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
					"3": testElement1{V: 3},
				},
				groundSetType: type1,
			},
			args:    args{e: testElement1{V: 4}},
			wantErr: false,
		},
		{
			name: "test2",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
					"3": testElement1{V: 3},
				},
				groundSetType: type1,
			},
			args:    args{e: testElement2{V: 4}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &Set{
				set:     tt.fields.set,
				setType: tt.fields.groundSetType,
			}
			if err := gs.Add(tt.args.e); (err != nil) != tt.wantErr {
				t.Errorf("Set.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	t.Run("test3", func(t *testing.T) {
		gs, _ := NewSet(type1, testElement1{V: 1}, testElement1{V: 2})
		if gs.Cardinality() != 2 {
			t.Errorf("cardinarity mismatch. expected:2, actural: %d", gs.Cardinality())
		}
		_ = gs.Add(testElement1{V: 1})
		if gs.Cardinality() != 2 {
			t.Errorf("cardinarity mismatch. expected:2, actural: %d", gs.Cardinality())
		}
		_ = gs.Add(testElement1{V: 3})
		if gs.Cardinality() != 3 {
			t.Errorf("cardinarity mismatch. expected:3, actural: %d", gs.Cardinality())
		}
	})
}

func TestSet_Clone(t *testing.T) {
	type fields struct {
		set           map[string]Element
		groundSetType ElementType
	}

	v1 := NewUnweightedVector([]float64{1, 1, 1, 1})
	v2 := NewUnweightedVector([]float64{0, 1, 1, 1})
	v3 := NewUnweightedVector([]float64{0, 0, 1, 1})
	tests := []struct {
		name   string
		fields fields
		want   *Set
	}{
		{
			name: "test1",
			fields: fields{
				set: map[string]Element{
					v1.Key(): v1,
					v2.Key(): v2,
					v3.Key(): v3,
				},
				groundSetType: VectorType,
			},
			want: &Set{
				set: map[string]Element{
					v1.Key(): v1,
					v2.Key(): v2,
					v3.Key(): v3,
				},
				setType: VectorType,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &Set{
				set:     tt.fields.set,
				setType: tt.fields.groundSetType,
			}
			if got := gs.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Set.Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Contains(t *testing.T) {
	type fields struct {
		set           map[string]Element
		groundSetType ElementType
	}
	type args struct {
		e []Element
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "test1",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
					"3": testElement1{V: 3},
				},
				groundSetType: type1,
			},
			args: args{
				[]Element{
					testElement1{V: 1},
					testElement1{V: 2},
					testElement1{V: 3},
				},
			},
			want: true,
		},
		{
			name: "test2",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
					"3": testElement1{V: 3},
				},
				groundSetType: type1,
			},
			args: args{
				[]Element{testElement1{V: 4}},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &Set{
				set:     tt.fields.set,
				setType: tt.fields.groundSetType,
			}
			if got := gs.Contains(tt.args.e...); got != tt.want {
				t.Errorf("Set.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Difference(t *testing.T) {
	type fields struct {
		set           map[string]Element
		groundSetType ElementType
	}
	type args struct {
		other *Set
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Set
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
					"3": testElement1{V: 3},
					"4": testElement1{V: 4},
				},
				groundSetType: type1,
			},
			args: args{
				other: &Set{
					set: map[string]Element{
						"1": testElement1{V: 1},
						"2": testElement1{V: 2},
					},
					setType: type1,
				},
			},
			want: &Set{
				set: map[string]Element{
					"3": testElement1{V: 3},
					"4": testElement1{V: 4},
				},
				setType: type1,
			},
			wantErr: false,
		},
		{
			name: "test2",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
				},
				groundSetType: type1,
			},
			args: args{
				other: &Set{
					set: map[string]Element{
						"1": testElement1{V: 1},
						"2": testElement1{V: 2},
						"3": testElement1{V: 3},
						"4": testElement1{V: 4},
					},
					setType: type1,
				},
			},
			want: &Set{
				set:     map[string]Element{},
				setType: type1,
			},
			wantErr: false,
		},
		{
			name: "test3",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
				},
				groundSetType: type1,
			},
			args: args{
				other: &Set{
					set: map[string]Element{
						"1": testElement1{V: 1},
						"2": testElement1{V: 2},
						"3": testElement1{V: 3},
						"4": testElement1{V: 4},
					},
					setType: type2,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &Set{
				set:     tt.fields.set,
				setType: tt.fields.groundSetType,
			}
			got, err := gs.Difference(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Set.Difference() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Set.Difference() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Equal(t *testing.T) {
	type fields struct {
		set           map[string]Element
		groundSetType ElementType
	}
	type args struct {
		other *Set
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "test1",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
					"3": testElement1{V: 3},
				},
				groundSetType: type1,
			},
			args: args{
				other: &Set{
					set: map[string]Element{
						"1": testElement1{V: 1},
						"2": testElement1{V: 2},
						"3": testElement1{V: 3},
					},
					setType: type1,
				},
			},
			want: true,
		},
		{
			name: "test2",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
					"3": testElement1{V: 3},
				},
				groundSetType: type1,
			},
			args: args{
				other: &Set{
					set: map[string]Element{
						"1": testElement1{V: 1},
						"2": testElement1{V: 2},
					},
					setType: type1,
				},
			},
			want: false,
		},
		{
			name: "test3",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
				},
				groundSetType: type1,
			},
			args: args{
				other: &Set{
					set: map[string]Element{
						"1": testElement1{V: 1},
						"2": testElement1{V: 2},
						"3": testElement1{V: 3},
					},
					setType: type1,
				},
			},
			want: false,
		},
		{
			name: "test4",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
					"3": testElement1{V: 3},
				},
				groundSetType: type1,
			},
			args: args{
				other: &Set{
					set: map[string]Element{
						"1": testElement2{V: 1},
						"2": testElement2{V: 2},
						"3": testElement2{V: 3},
					},
					setType: type2,
				},
			},
			// doesn't matter if setType are not equal
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &Set{
				set:     tt.fields.set,
				setType: tt.fields.groundSetType,
			}
			if got := gs.Equal(tt.args.other); got != tt.want {
				t.Errorf("Set.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Intersect(t *testing.T) {
	type fields struct {
		set           map[string]Element
		groundSetType ElementType
	}
	type args struct {
		other *Set
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Set
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
					"3": testElement1{V: 3},
					"4": testElement1{V: 4},
				},
				groundSetType: type1,
			},
			args: args{
				other: &Set{
					set: map[string]Element{
						"3": testElement1{V: 3},
						"4": testElement1{V: 4},
						"5": testElement1{V: 5},
						"6": testElement1{V: 6},
					},
					setType: type1,
				},
			},
			want: &Set{
				set: map[string]Element{
					"3": testElement1{V: 3},
					"4": testElement1{V: 4},
				},
				setType: type1,
			},
			wantErr: false,
		},
		{
			name: "test1",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
					"3": testElement1{V: 3},
					"4": testElement1{V: 4},
				},
				groundSetType: type1,
			},
			args: args{
				other: &Set{
					set: map[string]Element{
						"3": testElement1{V: 3},
						"4": testElement1{V: 4},
						"5": testElement1{V: 5},
						"6": testElement1{V: 6},
					},
					setType: type2,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &Set{
				set:     tt.fields.set,
				setType: tt.fields.groundSetType,
			}
			got, err := gs.Intersect(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Set.Intersect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Set.Intersect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_IsProperSubsetOf(t *testing.T) {
	type fields struct {
		set           map[string]Element
		groundSetType ElementType
	}
	type args struct {
		other *Set
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "test1",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
				},
				groundSetType: type1,
			},
			args: args{
				other: &Set{
					set: map[string]Element{
						"2": testElement1{V: 2},
						"3": testElement1{V: 3},
					},
					setType: type1,
				},
			},
			want: false,
		},
		{
			name: "test2",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
				},
				groundSetType: type1,
			},
			args: args{
				other: &Set{
					set: map[string]Element{
						"1": testElement1{V: 1},
						"2": testElement1{V: 2},
						"3": testElement1{V: 3},
					},
					setType: type1,
				},
			},
			want: true,
		},
		{
			name: "test3",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
					"3": testElement1{V: 3},
				},
				groundSetType: type1,
			},
			args: args{
				other: &Set{
					set: map[string]Element{
						"1": testElement1{V: 1},
						"2": testElement1{V: 2},
						"3": testElement1{V: 3},
					},
					setType: type1,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &Set{
				set:     tt.fields.set,
				setType: tt.fields.groundSetType,
			}
			if got := gs.IsProperSubsetOf(tt.args.other); got != tt.want {
				t.Errorf("Set.IsProperSubsetOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_IsProperSupersetOf(t *testing.T) {
	type fields struct {
		set           map[string]Element
		groundSetType ElementType
	}
	type args struct {
		other *Set
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "test1",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
				},
				groundSetType: type1,
			},
			args: args{
				other: &Set{
					set: map[string]Element{
						"2": testElement1{V: 2},
						"3": testElement1{V: 3},
					},
					setType: type1,
				},
			},
			want: false,
		},
		{
			name: "test2",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
					"3": testElement1{V: 3},
				},
				groundSetType: type1,
			},
			args: args{
				other: &Set{
					set: map[string]Element{
						"1": testElement1{V: 1},
						"2": testElement1{V: 2},
					},
					setType: type1,
				},
			},
			want: true,
		},
		{
			name: "test3",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
					"3": testElement1{V: 3},
				},
				groundSetType: type1,
			},
			args: args{
				other: &Set{
					set: map[string]Element{
						"1": testElement1{V: 1},
						"2": testElement1{V: 2},
						"3": testElement1{V: 3},
					},
					setType: type1,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &Set{
				set:     tt.fields.set,
				setType: tt.fields.groundSetType,
			}
			if got := gs.IsProperSupersetOf(tt.args.other); got != tt.want {
				t.Errorf("Set.IsProperSupersetOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_IsSubsetOf(t *testing.T) {
	type fields struct {
		set           map[string]Element
		groundSetType ElementType
	}
	type args struct {
		other *Set
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "test1",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
				},
				groundSetType: type1,
			},
			args: args{
				other: &Set{
					set: map[string]Element{
						"2": testElement1{V: 2},
						"3": testElement1{V: 3},
					},
					setType: type1,
				},
			},
			want: false,
		},
		{
			name: "test2",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
				},
				groundSetType: type1,
			},
			args: args{
				other: &Set{
					set: map[string]Element{
						"1": testElement1{V: 1},
						"2": testElement1{V: 2},
						"3": testElement1{V: 3},
					},
					setType: type1,
				},
			},
			want: true,
		},
		{
			name: "test3",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
					"3": testElement1{V: 3},
				},
				groundSetType: type1,
			},
			args: args{
				other: &Set{
					set: map[string]Element{
						"1": testElement1{V: 1},
						"2": testElement1{V: 2},
						"3": testElement1{V: 3},
					},
					setType: type1,
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &Set{
				set:     tt.fields.set,
				setType: tt.fields.groundSetType,
			}
			if got := gs.IsSubsetOf(tt.args.other); got != tt.want {
				t.Errorf("Set.IsSubsetOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_IsSupersetOf(t *testing.T) {
	type fields struct {
		set           map[string]Element
		groundSetType ElementType
	}
	type args struct {
		other *Set
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "test1",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
				},
				groundSetType: type1,
			},
			args: args{
				other: &Set{
					set: map[string]Element{
						"2": testElement1{V: 2},
						"3": testElement1{V: 3},
					},
					setType: type1,
				},
			},
			want: false,
		},
		{
			name: "test2",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
					"3": testElement1{V: 3},
				},
				groundSetType: type1,
			},
			args: args{
				other: &Set{
					set: map[string]Element{
						"1": testElement1{V: 1},
						"2": testElement1{V: 2},
					},
					setType: type1,
				},
			},
			want: true,
		},
		{
			name: "test3",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
					"3": testElement1{V: 3},
				},
				groundSetType: type1,
			},
			args: args{
				other: &Set{
					set: map[string]Element{
						"1": testElement1{V: 1},
						"2": testElement1{V: 2},
						"3": testElement1{V: 3},
					},
					setType: type1,
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &Set{
				set:     tt.fields.set,
				setType: tt.fields.groundSetType,
			}
			if got := gs.IsSupersetOf(tt.args.other); got != tt.want {
				t.Errorf("Set.IsSupersetOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Each(t *testing.T) {
	type fields struct {
		set           map[string]Element
		groundSetType ElementType
	}
	type args struct {
		f func(Element) bool
	}
	other := EmptySet(type2)
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Set
	}{
		{
			name: "test1",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
					"3": testElement1{V: 3},
				},
				groundSetType: type1,
			},
			args: args{
				func(e Element) bool {
					e0 := testElement2{V: e.Value().(int)}
					if e0.V > 1 {
						_ = other.Add(e0)
					}
					return true
				},
			},
			want: &Set{
				set: map[string]Element{
					"2": testElement2{V: 2},
					"3": testElement2{V: 3},
				},
				setType: type2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &Set{
				set:     tt.fields.set,
				setType: tt.fields.groundSetType,
			}
			gs.Each(tt.args.f)
			if !reflect.DeepEqual(other, tt.want) {
				t.Errorf("Set.Each(): %v, want %v", other, tt.want)
			}
		})
	}
}

func TestSet_Remove(t *testing.T) {
	type fields struct {
		set           map[string]Element
		groundSetType ElementType
	}
	tests := []struct {
		name   string
		fields fields
		args   []Element
		want   int
	}{
		{
			name: "test1",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
					"3": testElement1{V: 3},
				},
				groundSetType: type1,
			},
			args: []Element{testElement1{V: 1}, testElement1{V: 2}},
			want: 1,
		},
		{
			name: "test2",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
					"3": testElement1{V: 3},
				},
				groundSetType: type1,
			},
			args: []Element{testElement1{V: 3}, testElement1{V: 4}},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &Set{
				set:     tt.fields.set,
				setType: tt.fields.groundSetType,
			}
			for _, e := range tt.args {
				gs.Remove(e)
			}
			if gs.Cardinality() != tt.want {
				t.Errorf("Set.Remove(): %d, wantList %d", gs.Cardinality(), tt.want)
			}
		})
	}
}

func TestSet_SymmetricDifference(t *testing.T) {
	type fields struct {
		set           map[string]Element
		groundSetType ElementType
	}
	type args struct {
		other *Set
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Set
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
				},
				groundSetType: type1,
			},
			args: args{
				other: &Set{
					set: map[string]Element{
						"2": testElement1{V: 2},
						"3": testElement1{V: 3},
					},
					setType: type1,
				},
			},
			want: &Set{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"3": testElement1{V: 3},
				},
				setType: type1,
			},
			wantErr: false,
		},
		{
			name: "test2",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
				},
				groundSetType: type1,
			},
			args: args{
				other: &Set{
					set: map[string]Element{
						"2": testElement1{V: 2},
						"3": testElement1{V: 3},
					},
					setType: type2,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test3",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
					"3": testElement1{V: 3},
				},
				groundSetType: type1,
			},
			args: args{
				other: &Set{
					set: map[string]Element{
						"1": testElement1{V: 1},
						"2": testElement1{V: 2},
						"3": testElement1{V: 3},
					},
					setType: type1,
				},
			},
			want: &Set{
				set:     map[string]Element{},
				setType: type1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &Set{
				set:     tt.fields.set,
				setType: tt.fields.groundSetType,
			}
			got, err := gs.SymmetricDifference(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Set.SymmetricDifference() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Set.SymmetricDifference() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Union(t *testing.T) {
	type fields struct {
		set           map[string]Element
		groundSetType ElementType
	}
	type args struct {
		other *Set
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Set
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
				},
				groundSetType: type1,
			},
			args: args{
				other: &Set{
					set: map[string]Element{
						"2": testElement1{V: 2},
						"3": testElement1{V: 3},
					},
					setType: type1,
				},
			},
			want: &Set{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
					"3": testElement1{V: 3},
				},
				setType: type1,
			},
			wantErr: false,
		},
		{
			name: "test2",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
				},
				groundSetType: type1,
			},
			args: args{
				other: &Set{
					set: map[string]Element{
						"2": testElement1{V: 2},
						"3": testElement1{V: 3},
					},
					setType: type2,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test3",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
					"3": testElement1{V: 3},
				},
				groundSetType: type1,
			},
			args: args{
				other: &Set{
					set: map[string]Element{
						"1": testElement1{V: 1},
						"2": testElement1{V: 2},
						"3": testElement1{V: 3},
					},
					setType: type1,
				},
			},
			want: &Set{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
					"3": testElement1{V: 3},
				},
				setType: type1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &Set{
				set:     tt.fields.set,
				setType: tt.fields.groundSetType,
			}
			got, err := gs.Union(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Set.Union() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Set.Union() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_String(t *testing.T) {
	type fields struct {
		set           map[string]Element
		groundSetType ElementType
	}
	tests := []struct {
		name     string
		fields   fields
		wantList []string
	}{
		{
			name: "test1",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
					"3": testElement1{V: 3},
				},
				groundSetType: type1,
			},
			wantList: []string{
				"Set{\n 1\n 2\n 3\n}",
				"Set{\n 1\n 3\n 2\n}",
				"Set{\n 2\n 1\n 3\n}",
				"Set{\n 2\n 3\n 1\n}",
				"Set{\n 3\n 1\n 2\n}",
				"Set{\n 3\n 2\n 1\n}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &Set{
				set:     tt.fields.set,
				setType: tt.fields.groundSetType,
			}
			if got := gs.String(); !isStringIn(tt.wantList, got) {
				t.Errorf("Set.String() = %v, want one of %v", got, strings.Join(tt.wantList, "\n"))
			}
		})
	}
}

func TestSet_Pop(t *testing.T) {
	type fields struct {
		set           map[string]Element
		groundSetType ElementType
	}
	tests := []struct {
		name     string
		fields   fields
		wantList []Element
	}{
		{
			name: "test1",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
					"3": testElement1{V: 3},
				},
				groundSetType: type1,
			},
			wantList: []Element{
				testElement1{V: 1},
				testElement1{V: 2},
				testElement1{V: 3},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &Set{
				set:     tt.fields.set,
				setType: tt.fields.groundSetType,
			}
			if got := gs.Pop(); !isElementIn(tt.wantList, got) {
				t.Errorf("Set.Pop() = %v, want one of %v", got, tt.wantList)
			}
		})
	}
}

func TestSet_Choose(t *testing.T) {
	type fields struct {
		set           map[string]Element
		groundSetType ElementType
	}
	type args struct {
		f func(Element) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Element
	}{
		{
			name: "test1",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
					"3": testElement1{V: 3},
				},
				groundSetType: type1,
			},
			args: args{
				f: func(e Element) bool {
					return e.Value().(int) > 2
				},
			},
			want: testElement1{V: 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &Set{
				set:     tt.fields.set,
				setType: tt.fields.groundSetType,
			}
			if got := gs.Choose(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Set.Choose() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_CondSubset(t *testing.T) {
	type fields struct {
		set           map[string]Element
		groundSetType ElementType
	}
	type args struct {
		f func(Element) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Set
	}{
		{
			name: "test1",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{V: 1},
					"2": testElement1{V: 2},
					"3": testElement1{V: 3},
					"4": testElement1{V: 4},
					"5": testElement1{V: 5},
				},
				groundSetType: type1,
			},
			args: args{
				f: func(e Element) bool {
					return e.Value().(int) > 2
				},
			},
			want: &Set{
				set: map[string]Element{
					"3": testElement1{V: 3},
					"4": testElement1{V: 4},
					"5": testElement1{V: 5},
				},
				setType: type1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &Set{
				set:     tt.fields.set,
				setType: tt.fields.groundSetType,
			}
			if got := gs.CondSubset(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Set.CondSubset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func isStringIn(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func isElementIn(slice []Element, e Element) bool {
	for _, v := range slice {
		if reflect.DeepEqual(v, e) {
			return true
		}
	}
	return false
}
