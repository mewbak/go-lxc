/*
 * lxc_test.go: Go bindings for lxc
 *
 * Copyright © 2013, S.Çağlar Onur
 *
 * Authors:
 * S.Çağlar Onur <caglar@10ur.org>
 *
 * This library is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public
 * License as published by the Free Software Foundation; either
 * version 2.1 of the License, or (at your option) any later version.

 * This library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * Lesser General Public License for more details.

 * You should have received a copy of the GNU Lesser General Public
 * License along with this library; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301  USA
 */
package lxc

import (
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

const (
	CONTAINER_NAME               = "rubik"
	CLONE_CONTAINER_NAME         = "O"
	CLONE_OVERLAY_CONTAINER_NAME = "O_o"
	CONFIG_FILE_PATH             = "/var/lib/lxc"
	CONFIG_FILE_NAME             = "/var/lib/lxc/rubik/config"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func TestVersion(t *testing.T) {
	t.Logf("LXC version: %s", Version())
}

func TestDefaultConfigPath(t *testing.T) {
	if DefaultConfigPath() != CONFIG_FILE_PATH {
		t.Errorf("DefaultConfigPath failed...")
	}
}

func TestSetConfigPath(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)
	defer PutContainer(z)

	current_path := z.ConfigPath()
	z.SetConfigPath("/tmp")
	new_path := z.ConfigPath()

	if current_path == new_path {
		t.Errorf("SetConfigPath failed...")
	}
}

func TestConcurrentDefined_Negative(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i <= 100; i++ {
		wg.Add(1)
		go func() {
			z := NewContainer(strconv.Itoa(rand.Intn(10)))
			defer PutContainer(z)

			// sleep for a while to simulate some dummy work
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(250)))

			if z.Defined() {
				t.Errorf("Defined_Negative failed...")
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestDefined_Negative(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)
	defer PutContainer(z)

	if z.Defined() {
		t.Errorf("Defined_Negative failed...")
	}
}

func TestCreate(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)
	defer PutContainer(z)

	if !z.Create("ubuntu", "amd64", "quantal") {
		t.Errorf("Creating the container failed...")
	}
}

func TestClone(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)
	defer PutContainer(z)

	if !z.Clone(CLONE_CONTAINER_NAME, DIRECTORY) {
		t.Errorf("Cloning the DIRECTORY backed container failed...")
	}

	if !z.Clone(CLONE_OVERLAY_CONTAINER_NAME, OVERLAYFS) {
		t.Errorf("Cloning the OVERLAYFS backed container failed...")
	}
}

func TestConcurrentCreate(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			z := NewContainer(strconv.Itoa(i))
			defer PutContainer(z)

			// sleep for a while to simulate some dummy work
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(250)))

			if !z.Create("ubuntu", "amd64", "quantal") {
				t.Errorf("Creating the container (%d) failed...", i)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestContainerNames(t *testing.T) {
	t.Logf("Containers: %+v\n", ContainerNames())
}

func TestContainers(t *testing.T) {
	for _, v := range Containers() {
		t.Logf("%s: %s", v.Name(), v.State())
	}
}

func TestConcurrentStart(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			z := NewContainer(strconv.Itoa(i))
			defer PutContainer(z)

			z.SetDaemonize()
			z.Start(false)
			z.Wait(RUNNING, 30)
			if !z.Running() {
				t.Errorf("Starting the container failed...")
			}

			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestConfigFileName(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)
	defer PutContainer(z)
	if z.ConfigFileName() != CONFIG_FILE_NAME {
		t.Errorf("ConfigFileName failed...")
	}
}

func TestDefined_Positive(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)
	defer PutContainer(z)

	if !z.Defined() {
		t.Errorf("Defined_Positive failed...")
	}
}

func TestConcurrentDefined_Positive(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i <= 100; i++ {
		wg.Add(1)
		go func() {
			z := NewContainer(strconv.Itoa(rand.Intn(10)))
			defer PutContainer(z)

			// sleep for a while to simulate some dummy work
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(250)))

			if !z.Defined() {
				t.Errorf("Defined_Positive failed...")
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestInitPID_Negative(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)
	defer PutContainer(z)

	if z.InitPID() != -1 {
		t.Errorf("InitPID failed...")
	}
}

func TestStart(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)
	defer PutContainer(z)

	z.SetDaemonize()
	z.Start(false)

	z.Wait(RUNNING, 30)
	if !z.Running() {
		t.Errorf("Starting the container failed...")
	}
}

func TestSetDaemonize(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)
	defer PutContainer(z)

	z.SetDaemonize()
	if !z.Daemonize() {
		t.Errorf("Daemonize failed...")
	}
}

func TestInitPID_Positive(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)
	defer PutContainer(z)

	if z.InitPID() == -1 {
		t.Errorf("InitPID failed...")
	}
}

func TestName(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)
	defer PutContainer(z)

	if z.Name() != CONTAINER_NAME {
		t.Errorf("Name failed...")
	}
}

func TestFreeze(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)
	defer PutContainer(z)

	z.Freeze()

	z.Wait(FROZEN, 30)
	if z.State() != FROZEN {
		t.Errorf("Freezing the container failed...")
	}
}

func TestUnfreeze(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)
	defer PutContainer(z)

	z.Unfreeze()

	z.Wait(RUNNING, 30)
	if z.State() != RUNNING {
		t.Errorf("Unfreezing the container failed...")
	}
}

func TestLoadConfigFile(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)
	defer PutContainer(z)

	if !z.LoadConfigFile(CONFIG_FILE_NAME) {
		t.Errorf("LoadConfigFile failed...")
	}
}

func TestSaveConfigFile(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)
	defer PutContainer(z)

	if !z.SaveConfigFile(CONFIG_FILE_NAME) {
		t.Errorf("LoadConfigFile failed...")
	}
}

func TestConfigItem(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)
	defer PutContainer(z)

	if z.ConfigItem("lxc.utsname")[0] != CONTAINER_NAME {
		t.Errorf("ConfigItem failed...")
	}
}

func TestSetConfigItem(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)
	defer PutContainer(z)

	z.SetConfigItem("lxc.utsname", CONTAINER_NAME)
	if z.ConfigItem("lxc.utsname")[0] != CONTAINER_NAME {
		t.Errorf("ConfigItem failed...")
	}
}

func TestSetCgroupItem(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)
	defer PutContainer(z)

	max_mem := z.CgroupItem("memory.max_usage_in_bytes")[0]
	current_mem := z.CgroupItem("memory.limit_in_bytes")[0]
	z.SetCgroupItem("memory.limit_in_bytes", max_mem)
	new_mem := z.CgroupItem("memory.limit_in_bytes")[0]

	if new_mem == current_mem {
		t.Errorf("SetCgroupItem failed...")
	}
}

func TestClearConfigItem(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)
	defer PutContainer(z)

	z.ClearConfigItem("lxc.cap.drop")
	if z.ConfigItem("lxc.cap.drop")[0] != "" {
		t.Errorf("ClearConfigItem failed...")
	}
}

func TestKeys(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)
	defer PutContainer(z)

	keys := strings.Join(z.Keys("lxc.network.0"), " ")
	if !strings.Contains(keys, "mtu") {
		t.Errorf("Keys failed...")
	}
}

func TestNumberOfNetworkInterfaces(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)
	defer PutContainer(z)

	if z.NumberOfNetworkInterfaces() != 1 {
		t.Errorf("NumberOfNetworkInterfaces failed...")
	}
}

func TestMemoryUsageInBytes(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)
	defer PutContainer(z)

	mem_used, _ := z.MemoryUsageInBytes()
	swap_used, _ := z.SwapUsageInBytes()
	mem_limit, _ := z.MemoryLimitInBytes()
	swap_limit, _ := z.SwapLimitInBytes()

	t.Logf("Mem usage: %0.0f\n", mem_used)
	t.Logf("Mem usage: %s\n", mem_used)
	t.Logf("Swap usage: %0.0f\n", swap_used)
	t.Logf("Swap usage: %s\n", swap_used)
	t.Logf("Mem limit: %0.0f\n", mem_limit)
	t.Logf("Mem limit: %s\n", mem_limit)
	t.Logf("Swap limit: %0.0f\n", swap_limit)
	t.Logf("Swap limit: %s\n", swap_limit)

}

/*
func TestReboot(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)
	defer PutContainer(z)

	t.Logf("Rebooting the container...\n")
	z.Reboot()

	if z.Running() {
		t.Errorf("Rebooting the container failed...")
	}
}
*/

func TestConcurrentShutdown(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			z := NewContainer(strconv.Itoa(i))
			defer PutContainer(z)

			z.Shutdown(30)

			if z.Running() {
				t.Errorf("Shutting down the container failed...")
			}

			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestShutdown(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)
	defer PutContainer(z)

	z.Shutdown(30)

	if z.Running() {
		t.Errorf("Shutting down the container failed...")
	}
}

func TestStop(t *testing.T) {
	z := NewContainer(CONTAINER_NAME)
	defer PutContainer(z)

	z.Stop()

	if z.Running() {
		t.Errorf("Stopping the container failed...")
	}
}

func TestDestroy(t *testing.T) {
	z := NewContainer(CLONE_OVERLAY_CONTAINER_NAME)
	defer PutContainer(z)

	if !z.Destroy() {
		t.Errorf("Destroying the container failed...")
	}

	z = NewContainer(CLONE_CONTAINER_NAME)
	defer PutContainer(z)

	if !z.Destroy() {
		t.Errorf("Destroying the container failed...")
	}

	z = NewContainer(CONTAINER_NAME)
	defer PutContainer(z)

	if !z.Destroy() {
		t.Errorf("Destroying the container failed...")
	}
}

func TestConcurrentDestroy(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			z := NewContainer(strconv.Itoa(i))
			defer PutContainer(z)

			// sleep for a while to simulate some dummy work
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(250)))

			if !z.Destroy() {
				t.Errorf("Destroying the container failed...")
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
