
// starts resources to provision them.
build {
    from_sources = [
        "src.amazon-ebs.ubuntu-1604",
        "src.virtualbox-iso.ubuntu-1204",
    ]

    provision {
        communicator = "comm.ssh.vagrant"
       
        shell {

            string   = "string"
            int      = 42
            int64    = 43
            bool     = true
            trilean  = true
            duration = "10s"
            map_string_string {
                a = "b"
                c = "d"
            }
            slice_string = [
                "a",
                "b",
                "c",
            ]

            nested {
                string   = "string"
                int      = 42
                int64    = 43
                bool     = true
                trilean  = true
                duration = "10s"
                map_string_string {
                    a = "b"
                    c = "d"
                }
                slice_string = [
                    "a",
                    "b",
                    "c",
                ]
            }

            nested_slice {
            }
        }

        file {

            string   = "string"
            int      = 42
            int64    = 43
            bool     = true
            trilean  = true
            duration = "10s"
            map_string_string {
                a = "b"
                c = "d"
            }
            slice_string = [
                "a",
                "b",
                "c",
            ]

            nested {
                string   = "string"
                int      = 42
                int64    = 43
                bool     = true
                trilean  = true
                duration = "10s"
                map_string_string {
                    a = "b"
                    c = "d"
                }
                slice_string = [
                    "a",
                    "b",
                    "c",
                ]
            }

            nested_slice {
            }
        }
    }

    provision {
        communicator = "comm.ssh.secure"
    }

    post-process {
        amazon-import { 
            string   = "string"
            int      = 42
            int64    = 43
            bool     = true
            trilean  = true
            duration = "10s"
            map_string_string {
                a = "b"
                c = "d"
            }
            slice_string = [
                "a",
                "b",
                "c",
            ]

            nested {
                string   = "string"
                int      = 42
                int64    = 43
                bool     = true
                trilean  = true
                duration = "10s"
                map_string_string {
                    a = "b"
                    c = "d"
                }
                slice_string = [
                    "a",
                    "b",
                    "c",
                ]
            }

            nested_slice {
            }
        }
    }

}
