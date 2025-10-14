package browserfingerprint

import "fmt"

func GetChangeCanvasJavaScript(rgba [4]int) string {
	return fmt.Sprintf(`
            const getImageData = CanvasRenderingContext2D.prototype.getImageData;

            const noisify = (canvas, context) => {
                if (context) {
                    const shift = {
                        r: %d,
                        g: %d,
                        b: %d,
                        a: %d,
                    };
                    const width = canvas.width;
                    const height = canvas.height;

                    if (width && height) {
                        const imageData = getImageData.apply(context, [0, 0, width, height]);

                        for (let i = 0; i < height; i++) 
                            for (let j = 0; j < width; j++) {
                                const n = i * (width * 4) + j * 4;
                                imageData.data[n + 0] = imageData.data[n + 0] + shift.r;
                                imageData.data[n + 1] = imageData.data[n + 1] + shift.g;
                                imageData.data[n + 2] = imageData.data[n + 2] + shift.b;
                                imageData.data[n + 3] = imageData.data[n + 3] + shift.a;
                            }
                    

                        context.putImageData(imageData, 0, 0);
                    }
                }
            };

            HTMLCanvasElement.prototype.toBlob = new Proxy(HTMLCanvasElement.prototype.toBlob, {
                apply(target, self, args) {
                    noisify(self, self.getContext('2d'));

                    return Reflect.apply(target, self, args);
                },
            });

            HTMLCanvasElement.prototype.toDataURL = new Proxy(HTMLCanvasElement.prototype.toDataURL, {
                apply(target, self, args) {
                    noisify(self, self.getContext('2d'));

                    return Reflect.apply(target, self, args);
                },
            });

            CanvasRenderingContext2D.prototype.getImageData = new Proxy(
                CanvasRenderingContext2D.prototype.getImageData,
                {
                    apply(target, self, args) {
                        noisify(self.canvas, self);

                        return Reflect.apply(target, self, args);
                    },
                }
            );
`, rgba[0], rgba[1], rgba[2], rgba[3])
}
