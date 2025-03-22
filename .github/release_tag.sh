#!/bin/bash

# 检查是否提供了必要的参数
if [ $# -lt 2 ] || [ $# -gt 3 ]; then
    echo "用法: $0 <search_directory> <tag> [isLatest]"
    echo "例如: $0 ./projects v1.0.0"
    echo "或: $0 ./projects v1.0.0 latest"
    exit 1
fi

# 获取参数
search_dir=$1
tag=$2
is_latest=$3  # 第三个可选参数：是否为最新版本

# 检查搜索目录是否存在
if [ ! -d "$search_dir" ]; then
    echo "错误: 目录 '$search_dir' 不存在"
    exit 1
fi

# 遍历所有子目录
find . -type f -name "go.mod" | while read -r gomod_file; do
    # 获取go.mod文件所在的目录
    dir=$(dirname "$gomod_file")

    # 检查module name是否包含github.com/gone-io/goner
    if grep -q "module.*github.com/gone-io/goner" "$gomod_file"; then
        dir="${dir#.}"
        dir="${dir#/}"

        git_tag="$tag"
        if [ -n "$dir" ]; then
            git_tag="$dir/$tag"
        fi

        echo "git tag: $git_tag"
        git tag "$git_tag"
        if [ -n "$is_latest" ]; then
            latest_tag="latest"
            # 如果 dir 不为空
            if [ -n "$dir" ]; then
                latest_tag="$dir/latest"
            fi
            echo "git latest tag: $latest_tag"
            git tag "$latest_tag"
        fi
    fi
done

git push origin --tags