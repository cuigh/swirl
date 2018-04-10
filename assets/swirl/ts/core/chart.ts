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
          <a class="card-header-icon drag">
            <span class="icon">
              <i class="fas fa-bars has-text-grey-light" aria-hidden="true"></i>
            </span>
          </a>        
          <p class="card-header-title is-paddingless">${this.opts.title}</p>
          <a data-action="remove-chart" class="card-header-icon" aria-label="remove chart">
            <span class="icon">
              <i class="fas fa-trash has-text-danger" aria-hidden="true"></i>
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
                    splitNumber: 10,
                    splitLine: {show: false},
                },
                yAxis: {
                    type: 'value',
                    // boundaryGap: [0, '100%'],
                    splitLine: {
                        // show: false,
                        lineStyle: {
                            type: "dashed",
                        },
                    },
                    axisLabel: {
                        formatter: this.formatValue.bind(this),
                    },
                },
                color: this.getColors(),
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
                case "time:ns":
                    return value + 'ns';
                case "time:µs":
                    return value.toFixed(2) + 'µs';
                case "time:ms":
                    return value.toFixed(2) + 'ms';
                case "time:s":
                    if (value < 1) { // 1
                        return (value * 1000).toFixed(0) + 'ms';
                    } else {
                        return value.toFixed(2) + 's';
                    }
                case "time:m":
                    return value.toFixed(2) + 'm';
                case "time:h":
                    return value.toFixed(2) + 'h';
                case "time:d":
                    return value.toFixed(2) + 'd';
                case "size:bits":
                    value = value / 8; // fall-through
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
                case "size:kilobytes":
                    return value.toFixed(2) + 'K';
                case "size:megabytes":
                    return value.toFixed(2) + 'M';
                case "size:gigabytes":
                    return value.toFixed(2) + 'G';
                default:
                    return value.toFixed(2);
            }
        }

        private getColors(): string[] {
            // ECharts default colors
            // let colors = ['#c23531', '#2f4554', '#61a0a8', '#d48265', '#91c7ae', '#749f83', '#ca8622', '#bda29a', '#6e7074', '#546570', '#c4ccd3'];
            // let colors = [
            //     '#C1232B', '#B5C334', '#FCCE10', '#E87C25', '#27727B',
            //     '#FE8463', '#9BCA63', '#FAD860', '#F3A43B', '#60C0DD',
            //     '#D7504B', '#C6E579', '#F4E001', '#F0805A', '#26C0C0',
            // ];
            let colors = [
                '#45aaf2',
                '#6574cd',
                '#a55eea',
                '#f66d9b',
                '#cd201f',
                '#fd9644',
                '#f1c40f',
                '#7bd235',
                '#5eba00',
                '#2bcbba',
            ];
            this.shuffle(colors);
            return colors;
        }

        private shuffle(a: any) {
            let len = a.length;
            for (let i = 0; i < len - 1; i++) {
                let index = Math.floor(Math.random() * (len - i));
                let temp = a[index];
                a[index] = a[len - i - 1];
                a[len - i - 1] = temp;
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
                    center: ['40%', '50%'],
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
                    right: 20,
                    bottom: 30,
                },
                tooltip: {
                    formatter: (params: any) => {
                        let html = params[0].axisValueLabel + '<br/>';
                        for (let i = 0; i < params.length; i++) {
                            html += params[i].marker + params[i].seriesName + ': ' + this.formatValue(params[i].value[1]) + '<br/>';
                        }
                        return html;
                    },
                },
                yAxis: {
                    // min: 'dataMin',
                    max: 'dataMax',
                },
            });
        }

        setData(d: any): void {
            if (!d.series) {
                return;
            }

            d.series.forEach((s: any) => {
                s.type = this.opts.type;
                s.areaStyle = {
                    opacity: 0.3,
                };
                s.smooth = true;
                s.showSymbol = false;
                s.lineStyle = {
                    normal: {
                        width: 1,
                    }
                };
            });
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
        private readonly opts: ChartDashboardOptions;
        private charts: Chart[] = [];
        private timer: number;
        private dlg: ChartDialog;

        constructor(elem: string | Element | JQuery, charts: ChartOptions[], opts?: ChartDashboardOptions) {
            this.opts = $.extend(new ChartDashboardOptions(), opts);
            this.$panel = $(elem);
            this.dlg = new ChartDialog(this);

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
            let charts: any = [];
            this.$panel.children().each((index: number, elem: Element) => {
                let name = $(elem).data("name");
                for (let i = 0; i < this.charts.length; i++) {
                    let c = this.charts[i];
                    if (c.getOptions().name === name) {
                        charts.push({
                            name: c.getOptions().name,
                            width: c.getOptions().width,
                            height: c.getOptions().height,
                        });
                        break;
                    }
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

        getOptions(): ChartDashboardOptions {
            return this.opts;
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

    class ChartDialog {
        private dashboard: ChartDashboard;
        private fb: FilterBox;
        private charts: any;
        private $charts: JQuery;

        constructor(dashboard: ChartDashboard) {
            this.dashboard = dashboard;
            this.fb = new FilterBox("#txt-query", this.filterCharts.bind(this));
            $("#btn-add").click(this.showAddDlg.bind(this));
            $("#btn-add-chart").click(this.addChart.bind(this));
            $("#btn-save").click(() => {
                this.dashboard.save();
            });
        }

        private showAddDlg() {
            let $panel = $("#nav-charts");
            $panel.find("label.panel-block").remove();

            // load charts
            $ajax.get(`/system/chart/query`, {dashboard: this.dashboard.getOptions().name}).json((charts: any) => {
                for (let i = 0; i < charts.length; i++) {
                    let c = charts[i];
                    $panel.append(`<label class="panel-block">
          <input type="checkbox" value="${c.name}" data-index="${i}">${c.name}: ${c.title}
        </label>`);
                }
                this.charts = charts;
                this.$charts = $panel.find("label.panel-block");
            });

            let dlg = new Modal("#dlg-add-chart");
            dlg.show();
        }

        private filterCharts(text: string) {
            if (!text) {
                this.$charts.show();
                return;
            }

            this.$charts.each((i, elem) => {
                let $elem = $(elem);
                let texts: string[] = [
                    this.charts[i].name.toLowerCase(),
                    this.charts[i].title.toLowerCase(),
                    this.charts[i].desc.toLowerCase(),
                ];
                for (let i = 0; i < texts.length; i++) {
                    let index = texts[i].indexOf(text);
                    if (index >= 0) {
                        $elem.show();
                        return;
                    }
                }
                $elem.hide();
            })
        }

        private addChart() {
            this.$charts.each((i, e) => {
                if ($(e).find(":checked").length > 0) {
                    let c = this.charts[i];
                    this.dashboard.addGraph(c);
                }
            });
            Modal.close();
        }
    }
}