package matroid_intersection

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
)

const type1 ElementType = "TYPE1"

type testElement1 struct {
	Val    int
	Weight float64
}

func (e testElement1) GetType() ElementType {
	return type1
}

func (e testElement1) Key() string {
	return strconv.Itoa(e.Val)
}

func (e testElement1) Value() interface{} {
	return e.Val
}

const type2 ElementType = "TYPE2"

type testElement2 struct {
	Val    int
	Weight float64
}

func (e testElement2) GetType() ElementType {
	return type2
}

func (e testElement2) Key() string {
	return strconv.Itoa(e.Val)
}

func (e testElement2) Value() interface{} {
	return e.Val
}

func TestNewGroundSet(t *testing.T) {
	type args struct {
		t ElementType
		e []Element
	}
	tests := []struct {
		name    string
		args    args
		want    *GroundSet
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				t: type1,
				e: []Element{testElement1{Val: 1}, testElement1{Val: 2}, testElement1{Val: 3}},
			},
			want: &GroundSet{
				set: map[string]Element{
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
					"3": testElement1{Val: 3},
				},
				groundSetType: type1,
			},
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				t: type1,
				e: []Element{testElement1{Val: 1}, testElement1{Val: 2}, testElement2{Val: 3}},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGroundSet(tt.args.t, tt.args.e...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGroundSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGroundSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroundSet_Add(t *testing.T) {
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
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
					"3": testElement1{Val: 3},
				},
				groundSetType: type1,
			},
			args:    args{e: testElement1{Val: 4}},
			wantErr: false,
		},
		{
			name: "test2",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
					"3": testElement1{Val: 3},
				},
				groundSetType: type1,
			},
			args:    args{e: testElement2{Val: 4}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GroundSet{
				set:           tt.fields.set,
				groundSetType: tt.fields.groundSetType,
			}
			if err := gs.Add(tt.args.e); (err != nil) != tt.wantErr {
				t.Errorf("GroundSet.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	t.Run("test3", func(t *testing.T) {
		gs, _ := NewGroundSet(type1, testElement1{Val: 1}, testElement1{Val: 2})
		if gs.Cardinality() != 2 {
			t.Errorf("cardinarity mismatch. expected:2, actural: %d", gs.Cardinality())
		}
		_ = gs.Add(testElement1{Val: 1})
		if gs.Cardinality() != 2 {
			t.Errorf("cardinarity mismatch. expected:2, actural: %d", gs.Cardinality())
		}
		_ = gs.Add(testElement1{Val: 3})
		if gs.Cardinality() != 3 {
			t.Errorf("cardinarity mismatch. expected:3, actural: %d", gs.Cardinality())
		}
	})
}

func TestGroundSet_Clone(t *testing.T) {
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
		want   *GroundSet
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
			want: &GroundSet{
				set: map[string]Element{
					v1.Key(): v1,
					v2.Key(): v2,
					v3.Key(): v3,
				},
				groundSetType: VectorType,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GroundSet{
				set:           tt.fields.set,
				groundSetType: tt.fields.groundSetType,
			}
			if got := gs.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroundSet.Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroundSet_Contains(t *testing.T) {
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
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
					"3": testElement1{Val: 3},
				},
				groundSetType: type1,
			},
			args: args{
				[]Element{
					testElement1{Val: 1},
					testElement1{Val: 2},
					testElement1{Val: 3},
				},
			},
			want: true,
		},
		{
			name: "test2",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
					"3": testElement1{Val: 3},
				},
				groundSetType: type1,
			},
			args: args{
				[]Element{testElement1{Val: 4}},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GroundSet{
				set:           tt.fields.set,
				groundSetType: tt.fields.groundSetType,
			}
			if got := gs.Contains(tt.args.e...); got != tt.want {
				t.Errorf("GroundSet.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroundSet_Difference(t *testing.T) {
	type fields struct {
		set           map[string]Element
		groundSetType ElementType
	}
	type args struct {
		other *GroundSet
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *GroundSet
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
					"3": testElement1{Val: 3},
					"4": testElement1{Val: 4},
				},
				groundSetType: type1,
			},
			args: args{
				other: &GroundSet{
					set: map[string]Element{
						"1": testElement1{Val: 1},
						"2": testElement1{Val: 2},
					},
					groundSetType: type1,
				},
			},
			want: &GroundSet{
				set: map[string]Element{
					"3": testElement1{Val: 3},
					"4": testElement1{Val: 4},
				},
				groundSetType: type1,
			},
			wantErr: false,
		},
		{
			name: "test2",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
				},
				groundSetType: type1,
			},
			args: args{
				other: &GroundSet{
					set: map[string]Element{
						"1": testElement1{Val: 1},
						"2": testElement1{Val: 2},
						"3": testElement1{Val: 3},
						"4": testElement1{Val: 4},
					},
					groundSetType: type1,
				},
			},
			want: &GroundSet{
				set:           map[string]Element{},
				groundSetType: type1,
			},
			wantErr: false,
		},
		{
			name: "test3",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
				},
				groundSetType: type1,
			},
			args: args{
				other: &GroundSet{
					set: map[string]Element{
						"1": testElement1{Val: 1},
						"2": testElement1{Val: 2},
						"3": testElement1{Val: 3},
						"4": testElement1{Val: 4},
					},
					groundSetType: type2,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GroundSet{
				set:           tt.fields.set,
				groundSetType: tt.fields.groundSetType,
			}
			got, err := gs.Difference(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("GroundSet.Difference() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroundSet.Difference() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroundSet_Equal(t *testing.T) {
	type fields struct {
		set           map[string]Element
		groundSetType ElementType
	}
	type args struct {
		other *GroundSet
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
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
					"3": testElement1{Val: 3},
				},
				groundSetType: type1,
			},
			args: args{
				other: &GroundSet{
					set: map[string]Element{
						"1": testElement1{Val: 1},
						"2": testElement1{Val: 2},
						"3": testElement1{Val: 3},
					},
					groundSetType: type1,
				},
			},
			want: true,
		},
		{
			name: "test2",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
					"3": testElement1{Val: 3},
				},
				groundSetType: type1,
			},
			args: args{
				other: &GroundSet{
					set: map[string]Element{
						"1": testElement1{Val: 1},
						"2": testElement1{Val: 2},
					},
					groundSetType: type1,
				},
			},
			want: false,
		},
		{
			name: "test3",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
				},
				groundSetType: type1,
			},
			args: args{
				other: &GroundSet{
					set: map[string]Element{
						"1": testElement1{Val: 1},
						"2": testElement1{Val: 2},
						"3": testElement1{Val: 3},
					},
					groundSetType: type1,
				},
			},
			want: false,
		},
		{
			name: "test4",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
					"3": testElement1{Val: 3},
				},
				groundSetType: type1,
			},
			args: args{
				other: &GroundSet{
					set: map[string]Element{
						"1": testElement2{Val: 1},
						"2": testElement2{Val: 2},
						"3": testElement2{Val: 3},
					},
					groundSetType: type2,
				},
			},
			// doesn't matter if groundSetType are not equal
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GroundSet{
				set:           tt.fields.set,
				groundSetType: tt.fields.groundSetType,
			}
			if got := gs.Equal(tt.args.other); got != tt.want {
				t.Errorf("GroundSet.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroundSet_Intersect(t *testing.T) {
	type fields struct {
		set           map[string]Element
		groundSetType ElementType
	}
	type args struct {
		other *GroundSet
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *GroundSet
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
					"3": testElement1{Val: 3},
					"4": testElement1{Val: 4},
				},
				groundSetType: type1,
			},
			args: args{
				other: &GroundSet{
					set: map[string]Element{
						"3": testElement1{Val: 3},
						"4": testElement1{Val: 4},
						"5": testElement1{Val: 5},
						"6": testElement1{Val: 6},
					},
					groundSetType: type1,
				},
			},
			want: &GroundSet{
				set: map[string]Element{
					"3": testElement1{Val: 3},
					"4": testElement1{Val: 4},
				},
				groundSetType: type1,
			},
			wantErr: false,
		},
		{
			name: "test1",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
					"3": testElement1{Val: 3},
					"4": testElement1{Val: 4},
				},
				groundSetType: type1,
			},
			args: args{
				other: &GroundSet{
					set: map[string]Element{
						"3": testElement1{Val: 3},
						"4": testElement1{Val: 4},
						"5": testElement1{Val: 5},
						"6": testElement1{Val: 6},
					},
					groundSetType: type2,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GroundSet{
				set:           tt.fields.set,
				groundSetType: tt.fields.groundSetType,
			}
			got, err := gs.Intersect(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("GroundSet.Intersect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroundSet.Intersect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroundSet_IsProperSubset(t *testing.T) {
	type fields struct {
		set           map[string]Element
		groundSetType ElementType
	}
	type args struct {
		other *GroundSet
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
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
				},
				groundSetType: type1,
			},
			args: args{
				other: &GroundSet{
					set: map[string]Element{
						"2": testElement1{Val: 2},
						"3": testElement1{Val: 3},
					},
					groundSetType: type1,
				},
			},
			want: false,
		},
		{
			name: "test2",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
				},
				groundSetType: type1,
			},
			args: args{
				other: &GroundSet{
					set: map[string]Element{
						"1": testElement1{Val: 1},
						"2": testElement1{Val: 2},
						"3": testElement1{Val: 3},
					},
					groundSetType: type1,
				},
			},
			want: true,
		},
		{
			name: "test3",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
					"3": testElement1{Val: 3},
				},
				groundSetType: type1,
			},
			args: args{
				other: &GroundSet{
					set: map[string]Element{
						"1": testElement1{Val: 1},
						"2": testElement1{Val: 2},
						"3": testElement1{Val: 3},
					},
					groundSetType: type1,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GroundSet{
				set:           tt.fields.set,
				groundSetType: tt.fields.groundSetType,
			}
			if got := gs.IsProperSubset(tt.args.other); got != tt.want {
				t.Errorf("GroundSet.IsProperSubset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroundSet_IsProperSuperset(t *testing.T) {
	type fields struct {
		set           map[string]Element
		groundSetType ElementType
	}
	type args struct {
		other *GroundSet
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
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
				},
				groundSetType: type1,
			},
			args: args{
				other: &GroundSet{
					set: map[string]Element{
						"2": testElement1{Val: 2},
						"3": testElement1{Val: 3},
					},
					groundSetType: type1,
				},
			},
			want: false,
		},
		{
			name: "test2",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
					"3": testElement1{Val: 3},
				},
				groundSetType: type1,
			},
			args: args{
				other: &GroundSet{
					set: map[string]Element{
						"1": testElement1{Val: 1},
						"2": testElement1{Val: 2},
					},
					groundSetType: type1,
				},
			},
			want: true,
		},
		{
			name: "test3",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
					"3": testElement1{Val: 3},
				},
				groundSetType: type1,
			},
			args: args{
				other: &GroundSet{
					set: map[string]Element{
						"1": testElement1{Val: 1},
						"2": testElement1{Val: 2},
						"3": testElement1{Val: 3},
					},
					groundSetType: type1,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GroundSet{
				set:           tt.fields.set,
				groundSetType: tt.fields.groundSetType,
			}
			if got := gs.IsProperSuperset(tt.args.other); got != tt.want {
				t.Errorf("GroundSet.IsProperSuperset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroundSet_IsSubset(t *testing.T) {
	type fields struct {
		set           map[string]Element
		groundSetType ElementType
	}
	type args struct {
		other *GroundSet
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
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
				},
				groundSetType: type1,
			},
			args: args{
				other: &GroundSet{
					set: map[string]Element{
						"2": testElement1{Val: 2},
						"3": testElement1{Val: 3},
					},
					groundSetType: type1,
				},
			},
			want: false,
		},
		{
			name: "test2",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
				},
				groundSetType: type1,
			},
			args: args{
				other: &GroundSet{
					set: map[string]Element{
						"1": testElement1{Val: 1},
						"2": testElement1{Val: 2},
						"3": testElement1{Val: 3},
					},
					groundSetType: type1,
				},
			},
			want: true,
		},
		{
			name: "test3",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
					"3": testElement1{Val: 3},
				},
				groundSetType: type1,
			},
			args: args{
				other: &GroundSet{
					set: map[string]Element{
						"1": testElement1{Val: 1},
						"2": testElement1{Val: 2},
						"3": testElement1{Val: 3},
					},
					groundSetType: type1,
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GroundSet{
				set:           tt.fields.set,
				groundSetType: tt.fields.groundSetType,
			}
			if got := gs.IsSubset(tt.args.other); got != tt.want {
				t.Errorf("GroundSet.IsSubset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroundSet_IsSuperset(t *testing.T) {
	type fields struct {
		set           map[string]Element
		groundSetType ElementType
	}
	type args struct {
		other *GroundSet
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
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
				},
				groundSetType: type1,
			},
			args: args{
				other: &GroundSet{
					set: map[string]Element{
						"2": testElement1{Val: 2},
						"3": testElement1{Val: 3},
					},
					groundSetType: type1,
				},
			},
			want: false,
		},
		{
			name: "test2",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
					"3": testElement1{Val: 3},
				},
				groundSetType: type1,
			},
			args: args{
				other: &GroundSet{
					set: map[string]Element{
						"1": testElement1{Val: 1},
						"2": testElement1{Val: 2},
					},
					groundSetType: type1,
				},
			},
			want: true,
		},
		{
			name: "test3",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
					"3": testElement1{Val: 3},
				},
				groundSetType: type1,
			},
			args: args{
				other: &GroundSet{
					set: map[string]Element{
						"1": testElement1{Val: 1},
						"2": testElement1{Val: 2},
						"3": testElement1{Val: 3},
					},
					groundSetType: type1,
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GroundSet{
				set:           tt.fields.set,
				groundSetType: tt.fields.groundSetType,
			}
			if got := gs.IsSuperset(tt.args.other); got != tt.want {
				t.Errorf("GroundSet.IsSuperset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroundSet_Each(t *testing.T) {
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
		want   *GroundSet
	}{
		{
			name: "test1",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
					"3": testElement1{Val: 3},
				},
				groundSetType: type1,
			},
			args: args{
				func(e Element) bool {
					e0 := testElement2{Val: e.Value().(int)}
					if e0.Val > 1 {
						_ = other.Add(e0)
					}
					return true
				},
			},
			want: &GroundSet{
				set: map[string]Element{
					"2": testElement2{Val: 2},
					"3": testElement2{Val: 3},
				},
				groundSetType: type2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GroundSet{
				set:           tt.fields.set,
				groundSetType: tt.fields.groundSetType,
			}
			gs.Each(tt.args.f)
			if !reflect.DeepEqual(other, tt.want) {
				t.Errorf("GroundSet.Each(): %v, want %v", other, tt.want)
			}
		})
	}
}

func TestGroundSet_Remove(t *testing.T) {
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
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
					"3": testElement1{Val: 3},
				},
				groundSetType: type1,
			},
			args: []Element{testElement1{Val: 1}, testElement1{Val: 2}},
			want: 1,
		},
		{
			name: "test2",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
					"3": testElement1{Val: 3},
				},
				groundSetType: type1,
			},
			args: []Element{testElement1{Val: 3}, testElement1{Val: 4}},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GroundSet{
				set:           tt.fields.set,
				groundSetType: tt.fields.groundSetType,
			}
			for _, e := range tt.args {
				gs.Remove(e)
			}
			if gs.Cardinality() != tt.want {
				t.Errorf("GroundSet.Remove(): %d, wantList %d", gs.Cardinality(), tt.want)
			}
		})
	}
}

func TestGroundSet_SymmetricDifference(t *testing.T) {
	type fields struct {
		set           map[string]Element
		groundSetType ElementType
	}
	type args struct {
		other *GroundSet
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *GroundSet
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
				},
				groundSetType: type1,
			},
			args: args{
				other: &GroundSet{
					set: map[string]Element{
						"2": testElement1{Val: 2},
						"3": testElement1{Val: 3},
					},
					groundSetType: type1,
				},
			},
			want: &GroundSet{
				set: map[string]Element{
					"1": testElement1{Val: 1},
					"3": testElement1{Val: 3},
				},
				groundSetType: type1,
			},
			wantErr: false,
		},
		{
			name: "test2",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
				},
				groundSetType: type1,
			},
			args: args{
				other: &GroundSet{
					set: map[string]Element{
						"2": testElement1{Val: 2},
						"3": testElement1{Val: 3},
					},
					groundSetType: type2,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test3",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
					"3": testElement1{Val: 3},
				},
				groundSetType: type1,
			},
			args: args{
				other: &GroundSet{
					set: map[string]Element{
						"1": testElement1{Val: 1},
						"2": testElement1{Val: 2},
						"3": testElement1{Val: 3},
					},
					groundSetType: type1,
				},
			},
			want: &GroundSet{
				set:           map[string]Element{},
				groundSetType: type1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GroundSet{
				set:           tt.fields.set,
				groundSetType: tt.fields.groundSetType,
			}
			got, err := gs.SymmetricDifference(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("GroundSet.SymmetricDifference() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroundSet.SymmetricDifference() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroundSet_Union(t *testing.T) {
	type fields struct {
		set           map[string]Element
		groundSetType ElementType
	}
	type args struct {
		other *GroundSet
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *GroundSet
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
				},
				groundSetType: type1,
			},
			args: args{
				other: &GroundSet{
					set: map[string]Element{
						"2": testElement1{Val: 2},
						"3": testElement1{Val: 3},
					},
					groundSetType: type1,
				},
			},
			want: &GroundSet{
				set: map[string]Element{
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
					"3": testElement1{Val: 3},
				},
				groundSetType: type1,
			},
			wantErr: false,
		},
		{
			name: "test2",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
				},
				groundSetType: type1,
			},
			args: args{
				other: &GroundSet{
					set: map[string]Element{
						"2": testElement1{Val: 2},
						"3": testElement1{Val: 3},
					},
					groundSetType: type2,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test3",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
					"3": testElement1{Val: 3},
				},
				groundSetType: type1,
			},
			args: args{
				other: &GroundSet{
					set: map[string]Element{
						"1": testElement1{Val: 1},
						"2": testElement1{Val: 2},
						"3": testElement1{Val: 3},
					},
					groundSetType: type1,
				},
			},
			want: &GroundSet{
				set: map[string]Element{
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
					"3": testElement1{Val: 3},
				},
				groundSetType: type1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GroundSet{
				set:           tt.fields.set,
				groundSetType: tt.fields.groundSetType,
			}
			got, err := gs.Union(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("GroundSet.Union() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroundSet.Union() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroundSet_String(t *testing.T) {
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
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
					"3": testElement1{Val: 3},
				},
				groundSetType: type1,
			},
			wantList: []string{
				"GroundSet{\n 1\n 2\n 3\n}",
				"GroundSet{\n 1\n 3\n 2\n}",
				"GroundSet{\n 2\n 1\n 3\n}",
				"GroundSet{\n 2\n 3\n 1\n}",
				"GroundSet{\n 3\n 1\n 2\n}",
				"GroundSet{\n 3\n 2\n 1\n}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GroundSet{
				set:           tt.fields.set,
				groundSetType: tt.fields.groundSetType,
			}
			if got := gs.String(); !isStringIn(tt.wantList, got) {
				t.Errorf("GroundSet.String() = %v, want one of %v", got, strings.Join(tt.wantList, "\n"))
			}
		})
	}
}

func TestGroundSet_Pop(t *testing.T) {
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
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
					"3": testElement1{Val: 3},
				},
				groundSetType: type1,
			},
			wantList: []Element{
				testElement1{Val: 1},
				testElement1{Val: 2},
				testElement1{Val: 3},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GroundSet{
				set:           tt.fields.set,
				groundSetType: tt.fields.groundSetType,
			}
			if got := gs.Pop(); !isElementIn(tt.wantList, got) {
				t.Errorf("GroundSet.Pop() = %v, want one of %v", got, tt.wantList)
			}
		})
	}
}

func TestGroundSet_Choose(t *testing.T) {
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
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
					"3": testElement1{Val: 3},
				},
				groundSetType: type1,
			},
			args: args{
				f: func(e Element) bool {
					return e.Value().(int) > 2
				},
			},
			want: testElement1{Val: 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GroundSet{
				set:           tt.fields.set,
				groundSetType: tt.fields.groundSetType,
			}
			if got := gs.Choose(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroundSet.Choose() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroundSet_CondSubset(t *testing.T) {
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
		want   *GroundSet
	}{
		{
			name: "test1",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{Val: 1},
					"2": testElement1{Val: 2},
					"3": testElement1{Val: 3},
					"4": testElement1{Val: 4},
					"5": testElement1{Val: 5},
				},
				groundSetType: type1,
			},
			args: args{
				f: func(e Element) bool {
					return e.Value().(int) > 2
				},
			},
			want: &GroundSet{
				set: map[string]Element{
					"3": testElement1{Val: 3},
					"4": testElement1{Val: 4},
					"5": testElement1{Val: 5},
				},
				groundSetType: type1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GroundSet{
				set:           tt.fields.set,
				groundSetType: tt.fields.groundSetType,
			}
			if got := gs.CondSubset(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroundSet.CondSubset() = %v, want %v", got, tt.want)
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
