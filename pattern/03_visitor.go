package pattern

import "fmt"

/*
	Реализовать паттерн «посетитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Visitor_pattern
*/

// Паттерн Visitor - behavior pattern
// Заключается в добавлении нового поведения, через дополнительный класс, которому мы дали доступ.

// Класс Visitor не имеет доступ к приватным полям и методам, что
// менее гибко, чем метод на объекте.

// Visitable ---

type Node interface {
	accept(Visitor)
}

type File struct {
	name string
	size uint
}

func NewFile(name string, size uint) *File {
	return &File{name, size}
}

func (f *File) accept(v Visitor) {
	v.visitFile(f)
}

type Directory struct {
	name     string
	children []Node
}

func NewDir(name string, children ...Node) *Directory {
	return &Directory{name, children}
}

func (d *Directory) accept(v Visitor) {
	v.visitDirectory(d)
}

// Visitors ---

type Visitor interface {
	visitDirectory(*Directory)
	visitFile(*File)
}

type DiskUsage struct {
	sizeAcc  uint
	currPath string
	indent   uint
}

func (du *DiskUsage) HumanSize(size uint) string {
	idx := 0
	letter := []string{"B ", "KB", "MB", "GB"}
	for idx+1 < len(letter) && size > 1024 {
		size = size / 1024
		idx++
	}
	return fmt.Sprintf("%4d%s", size, letter[idx])
}

func (du *DiskUsage) visitFile(f *File) {
	du.sizeAcc += f.size
	fmt.Printf("%s %s%s\n", du.HumanSize(f.size), du.currPath, f.name)
}

func (du *DiskUsage) visitDirectory(d *Directory) {
	name := d.name
	sizeBackup := du.sizeAcc
	du.sizeAcc = 0
	defer func() { du.sizeAcc += sizeBackup }()

	pathBackup := du.currPath
	du.currPath += name + "/"
	defer func() { du.currPath = pathBackup }()

	for _, child := range d.children {
		child.accept(du)
	}

	fmt.Printf("%s %s\n", du.HumanSize(du.sizeAcc), du.currPath)
}

func useVisitor() {
	var root Node = NewDir("",
		NewDir("bin",
			NewFile("sh", 927242),
			NewFile("go", 14143932),
		),
		NewDir("home",
			NewDir("atlas",
				NewFile("dinosaur", 414255),
				NewFile("03997821.398", 4149484),
				NewDir("bmesa",
					NewFile("leak141928", 8535288),
				),
			),
		),
	)

	root.accept(&DiskUsage{})
}

// 905KB /bin/sh
//  13MB /bin/go
//  14MB /bin/
// 404KB /home/atlas/dinosaur
//   3MB /home/atlas/03997821.398
//   8MB /home/atlas/bmesa/leak141928
//   8MB /home/atlas/bmesa/
//  12MB /home/atlas/
//  12MB /home/
//  26MB /
