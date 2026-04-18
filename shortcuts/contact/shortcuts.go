package contact

import "codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"

func Shortcuts() []common.Shortcut {
	return []common.Shortcut{
		ContactSearchUser,
		ContactGetUser,
		ContactListDepartments,
		ContactGetDepartment,
		ContactListDepartmentUsers,
		ContactGetDepartmentUser,
		ContactListLectors,
	}
}
