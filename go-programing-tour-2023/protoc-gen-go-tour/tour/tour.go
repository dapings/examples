package tour

import (
	"github.com/golang/protobuf/protoc-gen-go/generator"
)

func init() {
	generator.RegisterPlugin(new(tour))
}

type tour struct {
	gen *generator.Generator
}

// 实现 generator.Plugin 接口，其中方法提供的FileDescriptor属性就包含了proto文件中的相应属性。
// FileDescriptor 属性所包含的信息，主要分为文件的描述信息、消息体的定义、枚举的定义、顶级扩展的定义、公开导入的文件中的所有类型的定义、注释、
// 导出的符号的完整列表(作为导出的对象到其符号的映射)、该文件包的导入路径、该文件的的Go软件包的名称，及是否为此文件生成proto3代码。

// Name 插件的名称
func (g *tour) Name() string {
	return "tour"
}

// Init 插件的初始化动作
func (g *tour) Init(gen *generator.Generator) {
	g.gen = gen
}

// Generate 生成文件所需的具体代码
func (g *tour) Generate(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}
}

// GenerateImports 生成文件所需的具体导入声明
func (g *tour) GenerateImports(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}
}
