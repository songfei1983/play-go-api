FROM amazoncorretto:17

# 设置时区为Asia/Tokyo
RUN ln -sf /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

# 安装必要的工具
RUN yum install -y wget tar gzip

# 安装 SBT
RUN wget https://github.com/sbt/sbt/releases/download/v1.9.8/sbt-1.9.8.tgz
RUN tar -xzvf sbt-1.9.8.tgz
RUN mv sbt /usr/local
ENV PATH="/usr/local/sbt/bin:${PATH}"

WORKDIR /gatling

# 复制 Gatling 项目文件
COPY build.sbt .
COPY project/ project/
COPY src/ src/

# 预下载依赖
RUN sbt update

# 设置入口点，添加 JVM 参数
ENTRYPOINT ["sbt", "-J--add-opens=java.base/java.lang=ALL-UNNAMED", "Gatling/test"]
