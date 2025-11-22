package types

const NGINX_CONF_TMP = `# Nginx 配置 - {{.SiteName}}
# 生成时间: {{.CreateTime}}

server {
    listen        {{.Port}};
    server_name   {{.Host}};
    root          "{{.FilePath}}";
    index         {{.Index}};

    # 字符编码
    charset utf-8;

    # 访问日志
    access_log  logs/{{.SiteName}}_access.log  combined;
    error_log   logs/{{.SiteName}}_error.log  warn;

    # 文件上传限制
    client_max_body_size 50m;

    # Gzip 压缩
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript application/x-javascript application/xml+rss application/json;

    # 静态文件缓存
    location ~* \.(jpg|jpeg|png|gif|ico|css|js|svg|woff|woff2|ttf|eot)$ {
        expires 7d;
        add_header Cache-Control "public, immutable";
    }

    # 主路由
    location / {
        try_files $uri $uri/ /{{.Index}};
        autoindex  off;
    }

    # 错误页面
    error_page 404 /404.html;
    error_page 500 502 503 504 /50x.html;
    location = /50x.html {
        root html;
    }

    # 禁止访问隐藏文件
    location ~ /\. {
        deny all;
        access_log off;
        log_not_found off;
    }

    # PHP 支持（如果需要）
    location ~ \.php$ {
        fastcgi_pass   127.0.0.1:9000;
        fastcgi_index  index.php;
        fastcgi_param  SCRIPT_FILENAME  $document_root$fastcgi_script_name;
        include        fastcgi_params;
    }

    # 禁止访问备份文件
    location ~* \.(bak|config|sql|fla|psd|ini|log|sh|inc|swp|dist)$ {
        deny all;
    }
}`
