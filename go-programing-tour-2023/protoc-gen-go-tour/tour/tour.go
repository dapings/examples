package tour

import (
	"github.com/golang/protobuf/protoc-gen-go/generator"
	pb "google.golang.org/protobuf/types/descriptorpb"
)

const (
	contextPkgPath = "context"
	errorsPkgPath  = "errors"
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

// Generate 生成文件所需的具体代码，是最核心的交通枢纽
func (g *tour) Generate(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}

	// 新增 errors 标准库的引入
	_ = g.gen.AddImport(errorsPkgPath)
	// string(g.gen.AddImport(contextPkgPath))
	// P方法的作用是将传入的参数，打印到所需生成的文件输出上，是一个相对原子的方法，并不会做过多的其他事情。
	// 从逻辑上看，方法主要是对proto文件的引入和版本信息进行了输出和定义，然后将最重要的服务(service)逻辑，通过FileDescriptorProto.Service循环调用generateService方法进行转换和输出。
	// generateService为主要的生成处理方法，其处理gRPC client生成的方法名为generateClientMethod。
}

// GenerateImports 生成文件所需的具体导入声明
func (g *tour) GenerateImports(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}
}

func (g *tour) generateOrgCodeMethod() {

}

func (g *tour) generateService(file *generator.FileDescriptor, service *pb.ServiceDescriptorProto, index int) {

	// 新增的 WithOrgCode 定义

	// 新增 OrgCode 定义
	g.generateOrgCodeMethod()
}

func (g *tour) generateClientMethod(servName, fullServName, serviceDescVar string, method *pb.MethodDescriptorProto, descExpr string) {
	// 新增对租户标识的获取和判断，若出现不存在或值为空的情况(也可针对实际业务场景自行定制)，则直接返回错误
	if !method.GetServerStreaming() && !method.GetClientStreaming() {

	}
}
