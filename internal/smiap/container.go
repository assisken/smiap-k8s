package smiap

type Container struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	DoNotRemove     bool   `json:"do_not_remove"`
	Cores           int    `json:"cores"`
	MemoryGb        int    `json:"memory_gb"`
	PartitionSizeGb int    `json:"partition_size_gb"`
}
