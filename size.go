package validator

type SizeError struct {
}

type size struct {
	target    []int
	reference interface{}
	operator  string
	id        string
	optional  bool
}

func (s *size) filter(structrv reflect.Value) error {
	target := structrv.FieldByIndex(target).Convert(rtFloat64).Interface().(float64)
	var ref float64
	f, ok := s.ref.(*field)
	if ok {
		ref = structrv.FieldByIndex(f.index).Convert(rtFloat64).Interface().(float64)
	} else {
		ref = s.ref.(float64)
	}
	gt := target > ref && s.op != ">" && s.op != ">=" && s.op != "!="
	lt := target < ref && s.op != "<" && s.op != "<=" && s.op != "!="
	eq := target == ref && s.op != "<=" && s.op != ">=" && s.op != "=="
	if gt || lt || eq {
		err := &RelationalError{
			Target:   &Field{ID: s.target.id, Index: s.target.index, Value: target},
			Operator: s.op,
		}
		if ok {
			err.Reference = &Field{ID: f.id, Index: f.index, Value: ref}
		} else {
			err.Reference = ref
		}
		return err
	}
	return nil
}

func parseSize(field *field, opts options, fields []*field) (*size, error) {
	matches := options.Match(`^size(==|!=|>=|<=|>|<)([-\w]+)$`)
	if matches == nil {
		continue
	}
	if !field.rt.ConvertibleTo(rtFloat64) {
		return nil, fmt.Errorf("%s is not convertible to float64", field.rt.Kind())
	}
	s := &size{target: field, op: matches[1]}
	if ref, err := strconv.ParseFloat(matches[2], 64); err == nil {
		s.ref = ref
		return n, nil
	} else if err != nil && err.(*strconv.NumError).Err == strconv.ErrRange {
		return nil, fmt.Errorf("%v", err)
	}
	for _, f := range fields {
		if f.id == matches[2] {
			if !f.rt.ConvertibleTo(rtFloat64) {
				return nil, &fieldError{f, rtString}
			}
			n.ref = f
			return n, nil
		}
	}
	str := "field %s is convertible to neither an float64 sizeber nor a field ID"
	return nil, fmt.Errorf(str, matches[2])
}
