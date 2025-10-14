package browserfingerprint

import (
	"fmt"
	"log"
	"time"
	_ "time/tzdata"
)

func GetChangeTimezoneJavaScript(timezone string, lang string, geo string) string {
	return fmt.Sprintf(`
		(() => {
			// 设置全局时区配置
			window.__rodTZConfig = {
				timeZone: '%s',
				language: '%s',
				geo: '%s'
			};
			
			// 覆盖 Date 对象
			const OriginalDate = window.Date;
			window.Date = function(...args) {
				if (new.target) {
					// 作为构造函数
					if (args.length === 0) {
						// 无参数调用时应用时区偏移
						const now = new OriginalDate();
						const targetOffset = -480; // %s 的 UTC 偏移（分钟）
						return new OriginalDate(now.getTime() + (targetOffset - now.getTimezoneOffset()) * 60000);
					}
					return new OriginalDate(...args);
				}
				// 作为函数调用
				return OriginalDate();
			};
			
			// 复制静态方法和属性
			Object.defineProperties(Date, Object.getOwnPropertyDescriptors(OriginalDate));
			
			// 覆盖 Date.now()
			const origDateNow = Date.now;
			Date.now = function() {
				return origDateNow() + (%d * 60000); // 添加时区偏移
			};
			
			// 覆盖 getTimezoneOffset
			const origGetTimezoneOffset = Date.prototype.getTimezoneOffset;
			Date.prototype.getTimezoneOffset = function() {
				return %d; // %s 的 UTC 偏移（分钟）
			};
			
			// 覆盖 Intl.DateTimeFormat
			const origDateTimeFormat = Intl.DateTimeFormat;
			Intl.DateTimeFormat = function(locales, options) {
				options = options || {};
				options.timeZone = options.timeZone || window.__rodTZConfig.timeZone;
				return new origDateTimeFormat(locales, options);
			};
			
			// 覆盖 navigator 属性
			Object.defineProperty(navigator, 'language', {
				get: () => window.__rodTZConfig.language,
				configurable: false
			});
			
			Object.defineProperty(navigator, 'languages', {
				get: () => [window.__rodTZConfig.language],
				configurable: false
			});
			
			// 覆盖控制台时区
			if (console._timeZone) {
				console._timeZone = window.__rodTZConfig.timeZone;
			}
		})();
	`,
		timezone,
		lang,
		geo,
		timezone,                      // 用于注释
		getUTCOffsetMinutes(timezone), // 时区偏移分钟数
		getUTCOffsetMinutes(timezone), // 时区偏移分钟数
		timezone)
}

// getUTCOffsetMinutes 获取时区的UTC偏移分钟数
func getUTCOffsetMinutes(timeZone string) int {
	// 使用时区名称获取实际偏移
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		log.Printf("无法加载时区 %s: %v, 使用默认值 -480", timeZone, err)
		return -480 // 默认太平洋时间
	}

	// 获取当前时间的偏移
	_, offset := time.Now().In(loc).Zone()
	return -offset / 60 // 转换为分钟并取反
}
