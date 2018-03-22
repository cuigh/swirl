namespace Swirl.Core {
    export class GraphOptions {
        type?: string = "line";
        unit?: string;
        width?: number = 12;
        height?: number = 50;
        colors?: string[];
    }

    export abstract class Graph {
        protected $elem: JQuery;
        protected opts?: GraphOptions;
        private readonly name: string;

        protected constructor(elem: string | Element | JQuery, opts?: GraphOptions) {
            this.$elem = $(elem);
            this.opts = $.extend(new GraphOptions(), opts);
            this.name = this.$elem.data("chart-name");
        }

        getName(): string {
            return this.name;
        }

        abstract setSize(w: number, h: number): void;

        abstract setData(d: any): void;
    }

    /**
     * Simple value
     */
    export class ValueGraph extends Graph {
        private $canvas: JQuery;

        constructor(elem: string | Element | JQuery, opts?: GraphOptions) {
            super(elem, opts);
        }

        setData(d: any): void {
        }

        setSize(w: number, h: number): void {
        }
    }

    export class ComplexGraph extends Graph {
        protected chart: Chart;
        protected ctx: CanvasRenderingContext2D;
        protected config: Chart.ChartConfiguration;
        protected static defaultColors = [
            'rgb(255, 99, 132)',    // red
            'rgb(75, 192, 192)',    // green
            'rgb(255, 159, 64)',    // orange
            'rgb(54, 162, 235)',    // blue
            'rgb(153, 102, 255)',   // purple
            'rgb(255, 205, 86)',    // yellow
            'rgb(201, 203, 207)',   // grey
        ];

        constructor(elem: string | Element | JQuery, opts?: GraphOptions) {
            super(elem, opts);

            if (!this.opts.colors) {
                this.opts.colors = ComplexGraph.defaultColors;
            }

            this.config = {
                type: opts.type,
                data: {},
                options: {
                    // title: {
                    //     // display: true,
                    //     text: opts.title || 'NONE'
                    // },
                    // legend: {
                    //     position: "bottom"
                    // },
                    animation: {
                        duration: 0,
                        // easing: 'easeOutBounce',
                    },
                    // tooltips: {
                    //     callbacks: {
                    //         label: function (tooltipItem: Chart.ChartTooltipItem, data: Chart.ChartData) {
                    //             let label = data.datasets[tooltipItem.datasetIndex].label || '';
                    //             if (label) {
                    //                 label += ': ';
                    //             }
                    //
                    //             let p = data.datasets[tooltipItem.datasetIndex].data[tooltipItem.index];
                    //             if (typeof p == "number") {
                    //                 label += ComplexGraph.formatValue(p, opts.unit);
                    //             } else {
                    //                 label += ComplexGraph.formatValue(p.y, opts.unit);
                    //             }
                    //             return label;
                    //         }
                    //     },
                    // }
                }
            };


            this.fillConfig();

            this.ctx = (<HTMLCanvasElement>$(elem).find("canvas").get(0)).getContext('2d');
            if (opts.height) {
                this.ctx.canvas.height = opts.height;
            }
            this.chart = new Chart(this.ctx, this.config);
        }

        setData(d: any): void {
        }

        setSize(w: number, h: number): void {
            this.ctx.canvas.height = h;
            this.chart.update();
        }

        protected fillConfig() {
        }

        protected static formatValue(value: number, unit: string): string {
            switch (unit) {
                case "percent:100":
                    return value.toFixed(1) + "%";
                case "percent:1":
                    return (value * 100).toFixed(1) + "%";
                case "time:s":
                    if (value < 1) { // 1
                        return (value * 1000).toFixed(0) + 'ms';
                    } else {
                        return value.toFixed(2) + 's';
                    }
                case "size:bytes":
                    if (value < 1024) { // 1K
                        return value.toString() + 'B';
                    } else if (value < 1048576) { // 1M
                        return (value / 1024).toFixed(2) + 'K';
                    } else if (value < 1073741824) { // 1G
                        return (value / 1048576).toFixed(2) + 'M';
                    } else {
                        return (value / 1073741824).toFixed(2) + 'G';
                    }
                default:
                    return value.toFixed(2);
            }
        }
    }

    /**
     * Pie etc.
     */
    export class VectorGraph extends ComplexGraph {
        protected fillConfig() {
            this.config.options.legend = {
                position: "right"
            };
            this.config.options.tooltips = {
                callbacks: {
                    label: (tooltipItem: Chart.ChartTooltipItem, data: Chart.ChartData): string => {
                        let label = data.labels[tooltipItem.index] + ": ";
                        let p = data.datasets[tooltipItem.datasetIndex].data[tooltipItem.index];
                        if (typeof p == "number") {
                            label += ComplexGraph.formatValue(p, this.opts.unit);
                        }
                        return label;
                    }
                },
            }
        }

        setData(d: any): void {
            // data = {
            //     datasets: [{
            //         data: [10, 20, 30]
            //     }],
            //
            //     // These labels appear in the legend and in the tooltips when hovering different arcs
            //     labels: [
            //         'Red',
            //         'Yellow',
            //         'Blue'
            //     ]
            // };
            this.config.data.datasets = [{
                data: d.data,
                backgroundColor: this.opts.colors,
            }];
            this.config.data.labels = d.labels;
            // this.config.data.datasets = [{
            //     data: [10, 20, 30],
            //     backgroundColor: this.opts.colors,
            // }];
            // this.config.data.labels = [
            //     'Red',
            //     'Yellow',
            //     'Blue'
            // ];
            this.chart.update();
        }
    }

    /**
     * Line/Bar etc.
     */
    export class MatrixGraph extends ComplexGraph {
        constructor(elem: string | Element | JQuery, opts?: GraphOptions) {
            super(elem, opts);
        }

        protected fillConfig() {
            this.config.options.scales = {
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
            };
            if (this.opts.unit) {
                this.config.options.scales.yAxes = [{
                    ticks: {
                        callback: (n: number) => ComplexGraph.formatValue(n, this.opts.unit),
                    }
                }];
                this.config.options.tooltips = {
                    callbacks: {
                        label: (tooltipItem: Chart.ChartTooltipItem, data: Chart.ChartData): string => {
                            let label = data.datasets[tooltipItem.datasetIndex].label + ": ";
                            let p: any = data.datasets[tooltipItem.datasetIndex].data[tooltipItem.index];
                            label += ComplexGraph.formatValue(p.y, this.opts.unit);
                            return label;
                        }
                    },
                }
            }
        }

        setData(d: any): void {
            let datasets = <Chart.ChartDataSets[]>(d);
            datasets.forEach((ds, i) => {
                let color = (i < this.opts.colors.length) ? this.opts.colors[i] : this.opts.colors[0];
                ds.backgroundColor = Chart.helpers.color(color).alpha(0.3).rgbString();
                ds.borderColor = color;
                ds.borderWidth = 2;
                ds.pointRadius = 2;
                // ds.fill = false;
            });
            this.config.data.datasets = d;
            this.chart.update();
        }
    }

    export class GraphFactory {
        static create(elem: string | Element | JQuery): Graph {
            let $elem = $(elem);
            let opts: GraphOptions = {
                type: $elem.data("chart-type"),
                unit: $elem.data("chart-unit"),
                height: $elem.data("chart-height"),
                colors: $elem.data("chart-colors"),
            };
            switch (opts.type) {
                case "value":
                    return new ValueGraph($elem, opts);
                case "line":
                case "bar":
                    return new MatrixGraph($elem, opts);
                case "pie":
                    return new VectorGraph($elem, opts);
            }
            return null;
        }
    }

    export class GraphPanelOptions {
        name: string;
        id?: string;
        time?: string = "30m";
        refreshInterval?: number = 15000;   // ms
    }

    export class GraphPanel {
        private opts: GraphPanelOptions;
        private charts: Graph[] = [];
        private timer: number;

        constructor(elems: string | Element | JQuery, opts?: GraphPanelOptions) {
            this.opts = $.extend(new GraphPanelOptions(), opts);

            $(elems).each((i, e) => {
                let g = GraphFactory.create(e);
                if (g != null) {
                    this.charts.push(g);
                }
            });

            this.refreshData();
        }

        private refreshData() {
            this.loadData();
            if (this.opts.refreshInterval > 0) {
                this.timer = setTimeout(this.refreshData.bind(this), this.opts.refreshInterval);
            }
        }

        refresh() {
            if (!this.timer) {
                this.loadData();
                if (this.opts.refreshInterval > 0) {
                    this.timer = setTimeout(this.refreshData.bind(this), this.opts.refreshInterval);
                }
            }
        }

        stop() {
            clearTimeout(this.timer);
            this.timer = 0;
        }

        setTime(time: string) {
            this.opts.time = time;
            this.loadData();
        }

        private loadData() {
            let args: any = {
                dashboard: this.opts.name,
                time: this.opts.time,
            };
            if (this.opts.id) {
                args.id = this.opts.id;
            }
            $ajax.get(`/system/chart/data`, args).json((d: { [index: string]: Chart.ChartDataSets[] }) => {
                $.each(this.charts, (i: number, g: Graph) => {
                    if (d[g.getName()]) {
                        g.setData(d[g.getName()]);
                    }
                });
            });
        }
    }
}