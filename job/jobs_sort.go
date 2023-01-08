package job

func (js Jobs) Len() int           { return len(js) }
func (js Jobs) Swap(i, j int)      { js[i], js[j] = js[j], js[i] }
func (js Jobs) Less(i, j int) bool { return js[i].estimated < js[j].estimated }
