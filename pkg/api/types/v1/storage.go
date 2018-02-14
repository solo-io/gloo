package v1

type StorableConfigObject interface {
	GetStorageRef() string
}

func (u *Upstream) GetStorageRef() string {
	return u.GetName()
}

func (v *VirtualHost) GetStorageRef() string {
	return v.GetName()
}
