package response

// Message排序：按照版本ID进行排序

type MessageSortByVersion []*MsgEntity

func (m MessageSortByVersion) Len() int {
	return len(m)
}

func (m MessageSortByVersion) Less(i, j int) bool {
	return m[i].VersionID < m[j].VersionID
}

func (m MessageSortByVersion) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

// Message排序：按照排序ID进行排序

type MessageSortBySortKey []*MsgEntity

func (m MessageSortBySortKey) Len() int {
	return len(m)
}

func (m MessageSortBySortKey) Less(i, j int) bool {
	return m[i].SortKey < m[j].SortKey
}

func (m MessageSortBySortKey) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

// Contact排序：按照版本ID进行排序

type ContactSortByVersion []*ContactEntity

func (c ContactSortByVersion) Len() int {
	return len(c)
}

func (c ContactSortByVersion) Less(i, j int) bool {
	return c[i].VersionID < c[j].VersionID
}

func (c ContactSortByVersion) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// Contact排序：按照排序ID进行排序

type ContactSortBySortKey []*ContactEntity

func (c ContactSortBySortKey) Len() int {
	return len(c)
}

func (c ContactSortBySortKey) Less(i, j int) bool {
	return c[i].SortKey < c[j].SortKey
}

func (c ContactSortBySortKey) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
