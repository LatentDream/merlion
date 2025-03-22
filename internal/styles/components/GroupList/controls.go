package grouplist

func (m *Model) prevGroup() {
	if m.selectedGroup > 0 {
		m.selectedGroup -= 1
		m.page = 0
		m.scrollOffset = 0
	}
}

func (m *Model) nextGroup() {
	if m.selectedGroup < len(m.Groups)-1 {
		m.selectedGroup += 1
		m.page = 0
		m.scrollOffset = 0
	}
}

func (m *Model) handleUpNavigation() {
	if m.selectedGroup >= 0 {
		if m.opennedGroup == nil {
			// No group open => move to previous group
			m.prevGroup()
		} else if m.selectedGroup == (*m.opennedGroup+1) && len(m.Groups[*m.opennedGroup].Items) == 0 {
			// Next group + no items in open => move to previous group
			m.prevGroup()
		} else if m.selectedGroup == (*m.opennedGroup+1) {
			// Next group + items in open group => move to last item
			m.prevGroup()
			lastItemIdx := len(m.Groups[*m.opennedGroup].Items) - 1
			m.selectedItem = &lastItemIdx
			m.ensureItemVisible()
		} else if m.selectedGroup == *m.opennedGroup && m.selectedItem == nil {
			// Current group + we are on the tag => move to previous group
			m.prevGroup()
		} else if m.selectedGroup == *m.opennedGroup {
			// Nav in the open group
			if len(m.Groups[*m.opennedGroup].Items) == 0 {
				// Safety check
				m.prevGroup()
			} else if m.selectedItem == nil {
				// Safety check
				m.selectedGroup = *m.opennedGroup
				lastItemIdx := len(m.Groups[*m.opennedGroup].Items) - 1
				m.selectedItem = &lastItemIdx
				m.ensureItemVisible()
			} else if *m.selectedItem == 0 {
				// First item => unselect item & stay on tag
				m.selectedItem = nil
			} else {
				// Prev Item
				*m.selectedItem -= 1
				m.ensureItemVisible()
			}
		} else {
			// Safety
			m.prevGroup()
		}
	}
}

func (m *Model) handleDownNavigation() {
	if m.selectedGroup <= len(m.Groups)-1 {
		if m.opennedGroup == nil {
			// No open group => Next
			m.nextGroup()
		} else if m.selectedGroup != *m.opennedGroup {
			// No in the open group => Next
			m.nextGroup()
		} else {
			if m.selectedItem == nil && len(m.Groups[m.selectedGroup].Items) == 0 {
				// No Item => Next
				m.nextGroup()
			} else if m.selectedItem == nil {
				// On the Group => First item
				var value int = 0
				m.selectedItem = &value
				m.ensureItemVisible()
			} else if *m.selectedItem == (len(m.Groups[m.selectedGroup].Items) - 1) {
				// Last item => Next group
				if m.selectedGroup != len(m.Groups)-1 {
					m.selectedItem = nil
					m.nextGroup()
				}
			} else {
				// Middle item => Next item
				*m.selectedItem += 1
				m.ensureItemVisible()
			}
		}
	}
}

func (m *Model) handleSelectItem() {
	if m.selectedItem == nil {
		if m.opennedGroup != nil && m.selectedGroup == *m.opennedGroup {
			m.opennedGroup = nil
		} else {
			m.opennedGroup = &m.selectedGroup
			m.selectedItem = nil
		}
		m.page = 0
		m.scrollOffset = 0
	}
}
