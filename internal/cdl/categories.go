package cdl

type Categories struct {
	Rooms                *parentCategory
	EmployeeIsRoomMember *parentCategory
	User                 *parentCategory
	FindMemberBy         *parentCategory
	ChatIDByMemberID     *parentCategory
	Message              *parentCategory
	Room                 *parentCategory
	MemberRole           *parentCategory
	UnitExistsByID       *parentCategory
	RoomExistsByID       *parentCategory
	MessageExists        *parentCategory
	Members              *parentCategory
}

func (d *Dataloader) ConfigureDataloader() {
	d.categories = &Categories{
		Rooms:                d.newRoomsCategory(),
		EmployeeIsRoomMember: d.newEmployeeIsRoomMemberCategory(),
		User:                 d.newUserCategory(),
		FindMemberBy:         d.newFindMemberByCategory(),
		ChatIDByMemberID:     d.newChatIDByMemberIDCategory(),
		Message:              d.newMessageCategory(),
		Room:                 d.newRoomCategory(),
		MemberRole:           d.newMemberRoleCategory(),
		UnitExistsByID:       d.newUnitExistsByIDCategory(),
		RoomExistsByID:       d.newRoomExistsByIDCategory(),
		MessageExists:        d.newMessageExistsCategory(),
		Members:              d.newMessageExistsCategory(),
	}
}

func (d *Dataloader) newMembersCategory() *parentCategory {
	c := d.newParentCategory()
	c.LoadFn = c.messageExists
	return c
}

func (d *Dataloader) newMessageExistsCategory() *parentCategory {
	c := d.newParentCategory()
	c.LoadFn = c.messageExists
	return c
}

func (d *Dataloader) newRoomExistsByIDCategory() *parentCategory {
	c := d.newParentCategory()
	c.LoadFn = c.roomExistsByID
	return c
}

func (d *Dataloader) newUnitExistsByIDCategory() *parentCategory {
	c := d.newParentCategory()
	c.LoadFn = c.unitExistsByID
	return c
}

func (d *Dataloader) newMemberRoleCategory() *parentCategory {
	c := d.newParentCategory()
	c.LoadFn = c.memberRole
	return c
}

func (d *Dataloader) newRoomCategory() *parentCategory {
	c := d.newParentCategory()
	c.LoadFn = c.room
	return c
}

func (d *Dataloader) newMessageCategory() *parentCategory {
	c := d.newParentCategory()
	c.LoadFn = c.message
	return c
}

func (d *Dataloader) newChatIDByMemberIDCategory() *parentCategory {
	c := d.newParentCategory()
	c.LoadFn = c.chatIDByMemberID
	return c
}

func (d *Dataloader) newFindMemberByCategory() *parentCategory {
	c := d.newParentCategory()
	c.LoadFn = c.findMemberBy
	return c
}

func (d *Dataloader) newUserCategory() *parentCategory {
	c := d.newParentCategory()
	c.LoadFn = c.user
	return c
}

func (d *Dataloader) newEmployeeIsRoomMemberCategory() *parentCategory {
	c := d.newParentCategory()
	c.LoadFn = c.employeeIsRoomMember
	return c
}

func (d *Dataloader) newRoomsCategory() *parentCategory {
	c := d.newParentCategory()
	c.LoadFn = c.rooms
	return c
}
