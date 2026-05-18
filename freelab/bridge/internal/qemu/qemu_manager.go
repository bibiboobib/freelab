package qemu

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// QemuManager управляет жизненным циклом виртуальных машин
type QemuManager struct {
	BaseImagesDir string // Папка, где лежат оригиналы (напр. mikrotik.qcow2, ubuntu.qcow2)
	OverlaysDir   string // Папка, куда складываются CoW диски для запущенных лаб
}

func NewQemuManager(baseDir, overlaysDir string) *QemuManager {
	// Создаем папки, если их нет
	os.MkdirAll(baseDir, 0755)
	os.MkdirAll(overlaysDir, 0755)
	
	return &QemuManager{
		BaseImagesDir: baseDir,
		OverlaysDir:   overlaysDir,
	}
}

// CreateOverlay создает Copy-On-Write клон (очень быстро, занимает КБ, а не ГБ)
func (qm *QemuManager) CreateOverlay(imageName string, vmID string) (string, error) {
	basePath := fmt.Sprintf("%s/%s", qm.BaseImagesDir, imageName)
	overlayPath := fmt.Sprintf("%s/%s.qcow2", qm.OverlaysDir, vmID)

	log.Printf("[QEMU] Создание CoW оверлея для VM %s на основе %s", vmID, imageName)

	// Команда: qemu-img create -f qcow2 -F qcow2 -b /base/image.qcow2 /overlays/vm1.qcow2
	cmd := exec.Command("qemu-img", "create", "-f", "qcow2", "-F", "qcow2", "-b", basePath, overlayPath)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ошибка qemu-img: %v, вывод: %s", err, string(output))
	}

	return overlayPath, nil
}

// StartVM запускает процесс виртуальной машины с привязкой к TAP интерфейсам
func (qm *QemuManager) StartVM(vmID string, overlayPath string, tapInterfaces []string) (*os.Process, error) {
	log.Printf("[QEMU] Запуск VM %s...", vmID)

	// Базовые аргументы (без графики, 512МБ ОЗУ, включаем KVM для скорости)
	args := []string{
		"-enable-kvm",
		"-m", "512M",
		"-nographic", // Запуск в фоне без окна QEMU
		"-drive", fmt.Sprintf("file=%s,format=qcow2,if=virtio", overlayPath),
	}

	// Динамически привязываем TAP-интерфейсы, которые мы создали через netlink
	for i, tap := range tapInterfaces {
		netdevArg := fmt.Sprintf("tap,id=net%d,ifname=%s,script=no,downscript=no", i, tap)
		deviceArg := fmt.Sprintf("virtio-net-pci,netdev=net%d,mac=52:54:00:12:34:%02x", i, i) // Генерируем MAC
		
		args = append(args, "-netdev", netdevArg)
		args = append(args, "-device", deviceArg)
	}

	cmd := exec.Command("qemu-system-x86_64", args...)
	
	// Важно: перенаправляем вывод, чтобы процесс не заблокировался
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Используем Start(), а не Run(), чтобы запустить процесс асинхронно в фоне и пойти дальше
	err := cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("не удалось запустить qemu-system-x86_64: %v", err)
	}

	log.Printf("[QEMU] VM %s успешно запущена, PID: %d", vmID, cmd.Process.Pid)
	return cmd.Process, nil
}
