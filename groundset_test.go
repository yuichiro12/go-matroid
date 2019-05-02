package matroid_intersection

import (
	"reflect"
	"strconv"
	"testing"
)

const type1 ElementType = "TYPE1"

type testElement1 struct {
	Value  int
	Weight float64
}

func (e testElement1) GetType() ElementType {
	return type1
}

func (e testElement1) Key() string {
	return strconv.Itoa(e.Value)
}

const type2 ElementType = "TYPE2"

type testElement2 struct {
	Value  int
	Weight float64
}

func (e testElement2) GetType() ElementType {
	return type2
}

func (e testElement2) Key() string {
	return strconv.Itoa(e.Value)
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
				e: []Element{testElement1{Value: 1}, testElement1{Value: 2}, testElement1{Value: 3}},
			},
			want: &GroundSet{
				set: map[string]Element{
					"1": testElement1{Value: 1},
					"2": testElement1{Value: 2},
					"3": testElement1{Value: 3},
				},
				groundSetType: type1,
			},
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				t: type1,
				e: []Element{testElement1{Value: 1}, testElement1{Value: 2}, testElement2{Value: 3}},
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
					"1": testElement1{Value: 1},
					"2": testElement1{Value: 2},
					"3": testElement1{Value: 3},
				},
				groundSetType: type1,
			},
			args:    args{e: testElement1{Value: 4}},
			wantErr: false,
		},
		{
			name: "test2",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{Value: 1},
					"2": testElement1{Value: 2},
					"3": testElement1{Value: 3},
				},
				groundSetType: type1,
			},
			args:    args{e: testElement2{Value: 4}},
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
		gs, _ := NewGroundSet(type1, testElement1{Value: 1}, testElement1{Value: 2})
		if gs.Cardinality() != 2 {
			t.Errorf("cardinarity mismatch. expected:2, actural: %d", gs.Cardinality())
		}
		_ = gs.Add(testElement1{Value: 1})
		if gs.Cardinality() != 2 {
			t.Errorf("cardinarity mismatch. expected:2, actural: %d", gs.Cardinality())
		}
		_ = gs.Add(testElement1{Value: 3})
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
					"1": testElement1{Value: 1},
					"2": testElement1{Value: 2},
					"3": testElement1{Value: 3},
				},
				groundSetType: type1,
			},
			args: args{
				[]Element{
					testElement1{Value: 1},
					testElement1{Value: 2},
					testElement1{Value: 3},
				},
			},
			want: true,
		},
		{
			name: "test2",
			fields: fields{
				set: map[string]Element{
					"1": testElement1{Value: 1},
					"2": testElement1{Value: 2},
					"3": testElement1{Value: 3},
				},
				groundSetType: type1,
			},
			args: args{
				[]Element{testElement1{Value: 4}},
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
		// TODO: Add test cases.
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
					"1": testElement1{Value: 1},
					"2": testElement1{Value: 2},
					"3": testElement1{Value: 3},
				},
				groundSetType: type1,
			},
			args: args{
				other: &GroundSet{
					set: map[string]Element{
						"1": testElement1{Value: 1},
						"2": testElement1{Value: 2},
						"3": testElement1{Value: 3},
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
					"1": testElement1{Value: 1},
					"2": testElement1{Value: 2},
					"3": testElement1{Value: 3},
				},
				groundSetType: type1,
			},
			args: args{
				other: &GroundSet{
					set: map[string]Element{
						"1": testElement1{Value: 1},
						"2": testElement1{Value: 2},
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
					"1": testElement1{Value: 1},
					"2": testElement1{Value: 2},
				},
				groundSetType: type1,
			},
			args: args{
				other: &GroundSet{
					set: map[string]Element{
						"1": testElement1{Value: 1},
						"2": testElement1{Value: 2},
						"3": testElement1{Value: 3},
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
					"1": testElement1{Value: 1},
					"2": testElement1{Value: 2},
					"3": testElement1{Value: 3},
				},
				groundSetType: type1,
			},
			args: args{
				other: &GroundSet{
					set: map[string]Element{
						"1": testElement2{Value: 1},
						"2": testElement2{Value: 2},
						"3": testElement2{Value: 3},
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GroundSet{
				set:           tt.fields.set,
				groundSetType: tt.fields.groundSetType,
			}
			gs.Each(tt.args.f)
		})
	}
}

func TestGroundSet_Iter(t *testing.T) {
	type fields struct {
		set           map[string]Element
		groundSetType ElementType
	}
	tests := []struct {
		name   string
		fields fields
		want   <-chan Element
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GroundSet{
				set:           tt.fields.set,
				groundSetType: tt.fields.groundSetType,
			}
			if got := gs.Iter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroundSet.Iter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroundSet_Remove(t *testing.T) {
	type fields struct {
		set           map[string]Element
		groundSetType ElementType
	}
	type args struct {
		e Element
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GroundSet{
				set:           tt.fields.set,
				groundSetType: tt.fields.groundSetType,
			}
			gs.Remove(tt.args.e)
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GroundSet{
				set:           tt.fields.set,
				groundSetType: tt.fields.groundSetType,
			}
			if got := gs.String(); got != tt.want {
				t.Errorf("GroundSet.String() = %v, want %v", got, tt.want)
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
		name   string
		fields fields
		want   Element
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GroundSet{
				set:           tt.fields.set,
				groundSetType: tt.fields.groundSetType,
			}
			if got := gs.Pop(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroundSet.Pop() = %v, want %v", got, tt.want)
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
