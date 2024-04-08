# go_common
简单常用的方法及注释整理


# 提供一个准确的版本号便于清晰引用--(使用tag)
```shell
tagVersion=v1.0.6

git add . 
git commit -m "update"
git tag ${tagVersion}
git push && git push origin ${tagVersion}


```