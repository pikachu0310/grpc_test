# Envoyの公式Dockerイメージを使用
FROM envoyproxy/envoy:v1.30.1

# ポート9090を開放
EXPOSE 9090

# envoy.yaml設定ファイルをコンテナにコピー
COPY envoy.yaml /etc/envoy/envoy.yaml

# Envoyを起動
CMD ["envoy", "-c", "/etc/envoy/envoy.yaml", "--log-level", "debug"]
