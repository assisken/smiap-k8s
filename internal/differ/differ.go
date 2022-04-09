package differ

func Diff(wantMap, haveMap map[string]DiffObject) (toCreate, toDelete, toUpdate []DiffItem, err error) {
	toCreate = []DiffItem{}
	toDelete = []DiffItem{}
	toUpdate = []DiffItem{}

	for id, want := range wantMap {
		have, ok := haveMap[id]

		wantItem, err := NewDiffItem(id, want)
		if err != nil {
			return nil, nil, nil, err
		}

		if !ok {
			wantItem.Reason = IsWanted
			toCreate = append(toCreate, wantItem)
			continue
		}

		haveItem, err := NewDiffItem(id, have)
		if err != nil {
			return nil, nil, nil, err
		}

		if wantItem.hash != haveItem.hash {
			wantItem.Reason = HashMismatch
			toUpdate = append(toUpdate, wantItem)
			delete(haveMap, id)
			continue
		}

		delete(haveMap, id)
	}

	for id, have := range haveMap {
		haveItem, err := NewDiffItem(id, have)
		if err != nil {
			return nil, nil, nil, err
		}

		haveItem.Reason = NotWanted
		toDelete = append(toDelete, haveItem)
	}

	return toCreate, toDelete, toUpdate, nil
}
