namespace Swirl.Core {
    export class ChartOptions {
        name: string;
        title: string;
        type?: string = "line";
        unit?: string;
        width?: number = 12;
        height?: number = 200;
        colors?: string[];
    }

    export abstract class Chart {
        protected $elem: JQuery;
        protected chart: echarts.ECharts;
        protected opts?: ChartOptions;

        protected constructor(opts: ChartOptions) {
            this.opts = $.extend(new ChartOptions(), opts);
            this.createElem();
        }

        private createElem() {
            this.$elem = $(`<div class="column is-${this.opts.width}" data-name="${this.opts.name}">
      <div class="card">
        <header class="card-header">
          <p class="card-header-title">${this.opts.title}</p>
          <a data-action="remove-chart" class="card-header-icon" aria-label="remove chart">
            <span class="icon">
              <i class="fas fa-times has-text-danger" aria-hidden="true"></i>
            </span>
          </a>
        </header>
        <div class="card-content" style="height: ${this.opts.height}px"></div>
      </div>
    </div>`);
        }

        init() {
            let opt = {
                legend: {
                    x: 'right',
                },
                tooltip: {
                    trigger: 'axis',
                    // formatter: function (params) {
                    //     params = params[0];
                    //     var date = new Date(params.name);
                    //     return date.getDate() + '/' + (date.getMonth() + 1) + '/' + date.getFullYear() + ' : ' + params.value[1];
                    // },
                    axisPointer: {
                        animation: false
                    }
                },
                xAxis: {
                    type: 'time',
                    splitLine: {show: false},
                },
                yAxis: {
                    type: 'value',
                    boundaryGap: [0, '100%'],
                    splitLine: {show: false},
                    axisLabel: {
                        formatter: this.formatValue.bind(this),
                    },
                },
            };
            this.config(opt);
            this.chart = echarts.init(<HTMLDivElement>this.$elem.find("div.card-content").get(0));
            this.chart.setOption(opt, true);
        }

        getElem(): JQuery {
            return this.$elem;
        }

        getOptions(): ChartOptions {
            return this.opts;
        }

        resize() {
            this.chart.resize();
        }

        abstract setData(d: any): void;

        protected config(opt: echarts.EChartOption) {
        }

        protected formatValue(value: number): string {
            switch (this.opts.unit) {
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
     * Gauge chart
     */
    export class GaugeChart extends Chart {
        constructor(opts: ChartOptions) {
            super(opts);
        }

        protected config(opt: echarts.EChartOption) {
            $.extend(true, opt, {
                grid: {
                    left: 0,
                    top: 20,
                    right: 0,
                    bottom: 0,
                },
                xAxis: [
                    {
                        show: false,
                    },
                ],
                yAxis: [
                    {
                        show: false,
                    },
                ],
            });
        }

        setData(d: any): void {
            this.chart.setOption({
                series: [
                    {
                        // name: d.name,
                        type: 'gauge',
                        radius: '100%',
                        center: ["50%", "58%"], // 仪表位置
                        max: d.value,
                        axisLabel: {show: false},
                        pointer: {show: false},
                        detail: {
                            offsetCenter: [0, 0],
                        },
                        data: [{value: d.value}]
                    }
                ]
            });
        }
    }

    /**
     * Pie chart.
     */
    export class VectorChart extends Chart {
        constructor(opts: ChartOptions) {
            super(opts);
        }

        protected config(opt: echarts.EChartOption) {
            $.extend(true, opt, {
                grid: {
                    left: 20,
                    top: 20,
                    right: 20,
                    bottom: 20,
                },
                legend: {
                    type: 'scroll',
                    orient: 'vertical',
                },
                tooltip: {
                    trigger: 'item',
                    formatter: (params: any, index: number): string => {
                        return params.name + ": " + this.formatValue(params.value);
                    },
                },
                xAxis: [
                    {
                        show: false,
                    },
                ],
                yAxis: [
                    {
                        show: false,
                    },
                ],
                series: [{
                    type: this.opts.type,
                    radius: '80%',
                    center: ['45%', '50%'],
                }],
            });
        }

        setData(d: any): void {
            this.chart.setOption({
                legend: {
                    data: d.legend,
                },
                series: [{
                    data: d.data,
                }],
            });
        }
    }

    /**
     * Line/Bar etc.
     */
    export class MatrixChart extends Chart {
        constructor(opts: ChartOptions) {
            super(opts);
        }

        protected config(opt: echarts.EChartOption) {
            $.extend(true, opt, {
                grid: {
                    left: 60,
                    top: 30,
                    right: 30,
                    bottom: 30,
                },
                tooltip: {
                    formatter: (params: any) => {
                        let html = params[0].axisValueLabel + '<br/>';
                        for (let i = 0; i < params.length; i++) {
                            html += params[i].marker + params[i].seriesName +': ' + this.formatValue(params[i].value[1]) + '<br/>';
                        }
                        return html;
                    },
                },
            });
        }

        setData(d: any): void {
            d.series.forEach((s: any) => s.type = this.opts.type);
            this.chart.setOption({
                legend: {
                    data: d.legend,
                },
                series: d.series,
            });
        }
    }

    export class ChartFactory {
        static create(opts: ChartOptions): Chart {
            switch (opts.type) {
                case "gauge":
                    return new GaugeChart(opts);
                case "line":
                case "bar":
                    return new MatrixChart(opts);
                case "pie":
                    return new VectorChart(opts);
            }
            return null;
        }
    }

    export class ChartDashboardOptions {
        name: string;
        key?: string;
        period?: number = 30;
        refreshInterval?: number = 15;   // ms
    }

    export class ChartDashboard {
        private $panel: JQuery;
        private opts: ChartDashboardOptions;
        private charts: Chart[] = [];
        private timer: number;

        constructor(elem: string | Element | JQuery, charts: ChartOptions[], opts?: ChartDashboardOptions) {
            this.opts = $.extend(new ChartDashboardOptions(), opts);
            this.$panel = $(elem);
            charts.forEach(opts => this.createGraph(opts));

            // events
            Dispatcher.bind(this.$panel).on("remove-chart", e => {
                let name = $(e.target).closest("div.column").data("name");
                Modal.confirm(`Are you sure to delete chart: <strong>${name}</strong>?`, "Delete chart", dlg => {
                    this.removeGraph(name);
                    dlg.close();
                });
            });
            $(window).resize(e => {
                $.each(this.charts, (i: number, g: Chart) => {
                    g.resize();
                });
            });

            this.refresh();
        }

        refresh() {
            if (!this.timer) {
                this.loadData();
                if (this.opts.refreshInterval > 0) {
                    this.timer = setTimeout(this.refreshData.bind(this), this.opts.refreshInterval * 1000);
                }
            }
        }

        private refreshData() {
            this.loadData();
            if (this.opts.refreshInterval > 0) {
                this.timer = setTimeout(this.refreshData.bind(this), this.opts.refreshInterval * 1000);
            }
        }

        stop() {
            clearTimeout(this.timer);
            this.timer = 0;
        }

        setPeriod(period: number) {
            this.opts.period = period;
            this.loadData();
        }

        addGraph(opts: ChartOptions) {
            this.createGraph(opts);
            this.loadData();
        }

        private createGraph(opts: ChartOptions) {
            for (let i = 0; i < this.charts.length; i++) {
                let chart = this.charts[i];
                if (chart.getOptions().name === opts.name) {
                    // chart already added.
                    return;
                }
            }

            let chart = ChartFactory.create(opts);
            if (chart != null) {
                this.$panel.append(chart.getElem());
                chart.init();
                this.charts.push(chart);
            }
        }

        removeGraph(name: string) {
            let index = -1;
            for (let i = 0; i < this.charts.length; i++) {
                let c = this.charts[i];
                if (c.getOptions().name === name) {
                    index = i;
                    break;
                }
            }

            if (index >= 0) {
                let $elem = this.charts[index].getElem();
                this.charts.splice(index, 1);
                $elem.remove();
            }
        }

        save() {
            let charts = this.charts.map(c => {
                return {
                    name: c.getOptions().name,
                    width: c.getOptions().width,
                    height: c.getOptions().height,
                }
            });
            let args = {
                name: this.opts.name,
                key: this.opts.key || '',
                charts: charts,
            };
            $ajax.post(`/system/chart/save_dashboard`, args).json<AjaxResult>((r: AjaxResult) => {
                if (r.success) {
                    Notification.show("success", "Successfully saved.");
                } else {
                    Notification.show("danger", r.message);
                }
            });
        }

        private loadData() {
            let args: any = {
                charts: this.charts.map(c => c.getOptions().name).join(","),
                period: this.opts.period,
            };
            if (this.opts.key) {
                args.key = this.opts.key;
            }
            $ajax.get(`/system/chart/data`, args).json((d: { [index: string]: any[] }) => {
                $.each(this.charts, (i: number, g: Chart) => {
                    if (d[g.getOptions().name]) {
                        g.setData(d[g.getOptions().name]);
                    }
                });
            });
        }
    }
}