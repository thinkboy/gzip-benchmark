server端提供三种压缩方式:

xxx.xxx.xxx.xxx:8080/oldGzip // golang自带

xxx.xxx.xxx.xxx:8080/klaGzip // 来源https://github.com/klauspost/compress/tree/master/gzip对gzip的优化

xxx.xxx.xxx.xxx:8080/myGzip // 使用对象池优化