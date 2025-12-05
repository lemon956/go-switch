#!/bin/bash
#
# go-switch uninit script
# 用于将 go-switch 环境完全回退到未初始化状态
#
# 功能:
# 1. 从所有支持的 shell 配置文件中移除 go-switch 的 source 命令
# 2. 删除 go-switch 的根目录（包括所有配置、环境文件、安装的 Go 版本等）
#

set -e

# 检查 root 权限
check_root() {
    if [ "$(id -u)" -ne 0 ]; then
        echo -e "\033[0;31m[ERROR]\033[0m 此脚本需要 root 权限运行"
        echo ""
        echo "请使用以下方式运行:"
        echo "  sudo $0 $*"
        echo ""
        exit 1
    fi
}

# 检查是否是帮助请求（不需要 root 权限）
is_help_request() {
    for arg in "$@"; do
        if [ "$arg" = "-h" ] || [ "$arg" = "--help" ]; then
            return 0
        fi
    done
    return 1
}

# 如果不是帮助请求，检查 root 权限
if ! is_help_request "$@"; then
    check_root "$@"
fi

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# go-switch 根目录
GO_SWITCH_DIR=".go-switch"
GO_SWITCH_ROOT="$HOME/$GO_SWITCH_DIR"

# 环境文件路径
GO_ENV_DIR="$GO_SWITCH_ROOT/environment"
SYSTEM_ENV_FILE="$GO_ENV_DIR/system"
SYSTEM_ENV_FILE_FISH="$GO_ENV_DIR/system.fish"

# shell 配置文件
BASHRC="$HOME/.bashrc"
BASH_PROFILE="$HOME/.bash_profile"
ZSHRC="$HOME/.zshrc"
FISH_CONFIG="$HOME/.config/fish/config.fish"

# 静默模式标志
SILENT=false
FORCE=false

# 打印带颜色的消息
print_info() {
    if [ "$SILENT" = false ]; then
        echo -e "${BLUE}[INFO]${NC} $1"
    fi
}

print_success() {
    if [ "$SILENT" = false ]; then
        echo -e "${GREEN}[SUCCESS]${NC} $1"
    fi
}

print_warning() {
    if [ "$SILENT" = false ]; then
        echo -e "${YELLOW}[WARNING]${NC} $1"
    fi
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
}

# 显示帮助信息
show_help() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "将 go-switch 环境回退到未初始化状态"
    echo ""
    echo "Options:"
    echo "  -y, --yes      跳过确认提示，直接执行"
    echo "  -s, --silent   静默模式，减少输出"
    echo "  -h, --help     显示此帮助信息"
    echo ""
    echo "此脚本将执行以下操作:"
    echo "  1. 从 shell 配置文件中移除 go-switch 的 source 命令"
    echo "     - ~/.bashrc (Linux)"
    echo "     - ~/.bash_profile (Mac)"
    echo "     - ~/.zshrc"
    echo "     - ~/.config/fish/config.fish"
    echo "  2. 删除 go-switch 根目录: $GO_SWITCH_ROOT"
    echo ""
}

# 解析命令行参数
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -y|--yes)
                FORCE=true
                shift
                ;;
            -s|--silent)
                SILENT=true
                shift
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                print_error "未知选项: $1"
                show_help
                exit 1
                ;;
        esac
    done
}

# 检查 go-switch 是否已初始化
check_initialized() {
    if [ ! -d "$GO_SWITCH_ROOT" ]; then
        print_warning "go-switch 目录不存在: $GO_SWITCH_ROOT"
        print_info "go-switch 可能尚未初始化，或已被删除"
        return 1
    fi
    
    if [ ! -f "$GO_SWITCH_ROOT/config/config.toml" ]; then
        print_warning "go-switch 配置文件不存在"
        print_info "go-switch 可能尚未完全初始化"
    fi
    
    return 0
}

# 从配置文件中移除包含指定模式的行
# 参数: $1 - 配置文件路径, $2 - 要匹配的模式
remove_line_from_config() {
    local config_file="$1"
    local pattern="$2"
    
    if [ ! -f "$config_file" ]; then
        print_info "配置文件不存在，跳过: $config_file"
        return 0
    fi
    
    # 检查文件中是否包含匹配的行
    if grep -q "$pattern" "$config_file" 2>/dev/null; then
        # 创建临时文件
        local temp_file=$(mktemp)
        
        # 过滤掉匹配的行
        grep -v "$pattern" "$config_file" > "$temp_file" || true
        
        # 替换原文件
        mv "$temp_file" "$config_file"
        
        print_success "已从 $config_file 中移除 go-switch 相关配置"
        return 0
    else
        print_info "配置文件中未发现 go-switch 相关配置: $config_file"
        return 0
    fi
}

# 清理 bash 配置
clean_bash_config() {
    print_info "清理 bash 配置..."
    
    # 匹配 source 命令 (bash/zsh 格式)
    # 例如: source /home/user/.go-switch/environment/system
    local pattern="source.*$GO_SWITCH_DIR/environment/system"
    
    # 清理 .bashrc
    remove_line_from_config "$BASHRC" "$pattern"
    
    # 清理 .bash_profile (Mac)
    remove_line_from_config "$BASH_PROFILE" "$pattern"
}

# 清理 zsh 配置
clean_zsh_config() {
    print_info "清理 zsh 配置..."
    
    # 匹配 source 命令
    local pattern="source.*$GO_SWITCH_DIR/environment/system"
    
    remove_line_from_config "$ZSHRC" "$pattern"
}

# 清理 fish 配置
clean_fish_config() {
    print_info "清理 fish 配置..."
    
    if [ ! -f "$FISH_CONFIG" ]; then
        print_info "fish 配置文件不存在，跳过: $FISH_CONFIG"
        return 0
    fi
    
    # fish 的 source 命令格式:
    # if test -f /home/user/.go-switch/environment/system.fish; source /home/user/.go-switch/environment/system.fish; end
    local pattern="$GO_SWITCH_DIR/environment/system.fish"
    
    if grep -q "$pattern" "$FISH_CONFIG" 2>/dev/null; then
        # 创建临时文件
        local temp_file=$(mktemp)
        
        # 过滤掉匹配的行
        grep -v "$pattern" "$FISH_CONFIG" > "$temp_file" || true
        
        # 替换原文件
        mv "$temp_file" "$FISH_CONFIG"
        
        print_success "已从 $FISH_CONFIG 中移除 go-switch 相关配置"
    else
        print_info "fish 配置文件中未发现 go-switch 相关配置"
    fi
}

# 清理所有 shell 配置
clean_all_shell_configs() {
    print_info "开始清理 shell 配置文件..."
    echo ""
    
    clean_bash_config
    clean_zsh_config
    clean_fish_config
    
    echo ""
    print_success "所有 shell 配置文件清理完成"
}

# 删除 go-switch 目录
delete_go_switch_dir() {
    print_info "删除 go-switch 目录: $GO_SWITCH_ROOT"
    
    if [ ! -d "$GO_SWITCH_ROOT" ]; then
        print_info "go-switch 目录不存在，无需删除"
        return 0
    fi
    
    # 显示将要删除的内容
    if [ "$SILENT" = false ]; then
        echo ""
        print_info "将要删除以下内容:"
        
        if [ -d "$GO_SWITCH_ROOT/config" ]; then
            echo "  - config/          (配置文件)"
        fi
        if [ -d "$GO_SWITCH_ROOT/environment" ]; then
            echo "  - environment/     (环境变量文件)"
        fi
        if [ -d "$GO_SWITCH_ROOT/gos" ]; then
            local go_versions=$(ls -1 "$GO_SWITCH_ROOT/gos" 2>/dev/null | wc -l)
            echo "  - gos/             (已安装的 Go 版本: $go_versions 个)"
        fi
        if [ -d "$GO_SWITCH_ROOT/go" ]; then
            echo "  - go/              (GOPATH 目录)"
        fi
        if [ -L "$GO_SWITCH_ROOT/current" ] || [ -d "$GO_SWITCH_ROOT/current" ]; then
            echo "  - current          (当前版本软链接)"
        fi
        echo ""
    fi
    
    # 执行删除
    rm -rf "$GO_SWITCH_ROOT"
    
    if [ ! -d "$GO_SWITCH_ROOT" ]; then
        print_success "go-switch 目录已删除"
    else
        print_error "删除 go-switch 目录失败"
        return 1
    fi
}

# 确认操作
confirm_action() {
    if [ "$FORCE" = true ]; then
        return 0
    fi
    
    echo ""
    print_warning "此操作将完全删除 go-switch 环境，包括:"
    echo "  - 所有已安装的 Go 版本"
    echo "  - go-switch 配置文件"
    echo "  - GOPATH 目录 (如果由 go-switch 管理)"
    echo "  - shell 配置文件中的 go-switch 相关配置"
    echo ""
    
    read -p "确定要继续吗? [y/N] " -n 1 -r
    echo ""
    
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_info "操作已取消"
        exit 0
    fi
}

# 显示完成信息
show_completion_message() {
    echo ""
    echo "=============================================="
    print_success "go-switch 环境已完全回退!"
    echo "=============================================="
    echo ""
    print_info "请执行以下操作以使更改生效:"
    echo ""
    echo "  1. 重新打开终端，或执行以下命令刷新配置:"
    echo ""
    echo "     Bash:  source ~/.bashrc"
    echo "     Zsh:   source ~/.zshrc"
    echo "     Fish:  source ~/.config/fish/config.fish"
    echo ""
    echo "  2. 如果您需要使用 Go，请重新安装官方 Go 或重新初始化 go-switch"
    echo ""
}

# 主函数
main() {
    parse_args "$@"
    
    echo ""
    echo "=============================================="
    echo "       go-switch 环境回退脚本"
    echo "=============================================="
    echo ""
    
    # 检查是否初始化
    if ! check_initialized; then
        # 即使目录不存在，也尝试清理 shell 配置
        print_info "尝试清理可能残留的 shell 配置..."
        clean_all_shell_configs
        print_info "清理完成"
        exit 0
    fi
    
    # 确认操作
    confirm_action
    
    echo ""
    
    # 清理 shell 配置
    clean_all_shell_configs
    
    echo ""
    
    # 删除 go-switch 目录
    delete_go_switch_dir
    
    # 显示完成信息
    show_completion_message
}

# 运行主函数
main "$@"

