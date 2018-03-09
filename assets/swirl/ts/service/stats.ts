///<reference path="../core/core.ts" />
namespace Swirl.Service {
    class MetricChartOptions {
        type?: string = "line";
        height?: number = 50;
        title?: string;
        labelX?: string;
        labelY?: string;
        tickY?: (value: number) => string;
    }

    class MetricData {
        cpu?: Chart.ChartDataSets[];
        memory?: Chart.ChartDataSets[];
    }

    class MetricChart {
        private chart: any;
        private readonly config: any;
        private colors = [
            'rgb(255, 99, 132)',    // red
            'rgb(75, 192, 192)',    // green
            'rgb(255, 159, 64)',    // orange
            'rgb(54, 162, 235)',    // blue
            'rgb(153, 102, 255)',   // purple
            'rgb(255, 205, 86)',    // yellow
            'rgb(201, 203, 207)',   // grey
        ];

        constructor(elem: string | Element, opts: MetricChartOptions) {
            opts = $.extend(new MetricChartOptions(), opts);
            this.config = {
                type: opts.type,
                data: {},
                options: {
                    title: {
                        // display: true,
                        text: opts.title || 'NONE'
                    },
                    // legend: {
                    //     position: "bottom"
                    // },
                    animation: {
                        duration: 0,
                        // easing: 'easeOutBounce',
                    },
                    scales: {
                        xAxes: [{
                            type: 'time',
                            time: {
                                // min: new Date(),
                                // max: new Date(),
                                unit: 'minute',
                                tooltipFormat: 'YYYY/MM/DD HH:mm:ss',
                                displayFormats: {
                                    minute: 'HH:mm'
                                }
                            },
                        }],
                        yAxes: [{}]
                    },
                }
            };
            if (opts.labelX) {
                this.config.options.scales.xAxes[0].scaleLabel = {
                    display: true,
                    labelString: opts.labelX,
                }
            }
            if (opts.labelY) {
                this.config.options.scales.yAxes[0].scaleLabel = {
                    display: true,
                    labelString: opts.labelY,
                }
            }
            if (opts.tickY) {
                this.config.options.scales.yAxes[0].ticks = {
                    callback: opts.tickY,
                }
            }

            let ctx = (<HTMLCanvasElement>$(elem).get(0)).getContext('2d');
            if (opts.height) {
                ctx.canvas.height = opts.height;
            }
            this.chart = new Chart(ctx, this.config);
        }

        public setData(datasets: Chart.ChartDataSets[]) {
            datasets.forEach((ds, i) => {
                let color = (i < this.colors.length) ? this.colors[i] : this.colors[0];
                ds.backgroundColor = Chart.helpers.color(color).alpha(0.3).rgbString();
                ds.borderColor = color;
                ds.borderWidth = 2;
                ds.pointRadius = 2;
                // ds.fill = false;
            });
            this.config.data.datasets = datasets;
            this.chart.update();
        }
    }

    export class StatsPage {
        private cpu: MetricChart;
        private memory: MetricChart;
        private timer: number;

        constructor() {
            let $cb_time = $("#cb-time");
            if ($cb_time.length == 0) {
                return;
            }

            $cb_time.change(this.loadData.bind(this));
            $("#cb-refresh").change(() => {
                if (this.timer) {
                    clearTimeout(this.timer);
                    this.timer = null;
                } else {
                    this.refreshData();
                }
            });

            this.cpu = new MetricChart("#canvas-cpu", {
                tickY: function (value: number): string {
                    return value + '%';
                },
            });
            this.memory = new MetricChart("#canvas-memory", {
                tickY: function (value: number): string {
                    return value < 1024 ? (value + 'M') : (value / 1024) + 'G';
                },
            });
            this.refreshData();
        }

        private refreshData() {
            this.loadData();
            if ($("#cb-refresh").prop("checked")) {
                this.timer = setTimeout(this.refreshData.bind(this), 15000);
            }
        }

        private loadData() {
            let time = $("#cb-time").val();
            $ajax.get(`metrics`, {time: time}).json((d: MetricData) => {
                if (d.cpu) {
                    this.cpu.setData(d.cpu);
                }
                if (d.memory) {
                    this.memory.setData(d.memory);
                }
            });
        }
    }
}