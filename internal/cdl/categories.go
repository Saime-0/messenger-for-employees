package cdl

type Categories struct {
	Rooms                *parentCategory
	EmployeeIsRoomMember *parentCategory
	Employee             *parentCategory
	Message              *parentCategory
	Room                 *parentCategory
	RoomExistsByID       *parentCategory
	MessageExists        *parentCategory
	Members              *parentCategory
	Tags                 *parentCategory
}

func (d *Dataloader) ConfigureDataloader() {
	d.categories = &Categories{
		Rooms:                d.newRoomsCategory(),
		EmployeeIsRoomMember: d.newEmployeeIsRoomMemberCategory(),
		Employee:             d.newEmployeeCategory(),
		Message:              d.newMessageCategory(),
		Room:                 d.newRoomCategory(),
		RoomExistsByID:       d.newRoomExistsByIDCategory(),
		MessageExists:        d.newMessageExistsCategory(),
		Members:              d.newMessageExistsCategory(),
		Tags:                 d.newTagsCategory(),
	}
}

func (d *Dataloader) newTagsCategory() *parentCategory {
	c := d.newParentCategory()
	c.LoadFn = c.tags
	return c
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

func (d *Dataloader) newEmployeeCategory() *parentCategory {
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
