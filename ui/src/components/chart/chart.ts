import * as echarts from "echarts";
import { ChartInfo } from "@/api/chart";
import { store } from "@/store";

export abstract class Chart {
  protected info: ChartInfo;
  protected chart: echarts.ECharts;

  protected constructor(dom: HTMLElement, info: ChartInfo) {
    this.info = info;
    this.chart = echarts.init(dom, store.state.preference.theme || 'light');
    this.setOption()
  }

  private setOption() {
    const basicOpt: echarts.EChartsOption = {
      legend: {
        align: 'right',
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
        splitLine: { show: false },
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
    this.chart.setOption(basicOpt, true, true);
    this.chart.setOption(this.config(), false, false);
  }

  private getColors() {
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
    colors.sort((a: string, b: string) => Math.random() - 0.5)
    return colors
  }

  abstract setData(d: any): void;

  protected abstract config(): echarts.EChartsOption

  protected formatValue(value: number): string {
    if (value < 0) {
      return "-" + this.formatValue(-value);
    }

    switch (this.info.unit) {
      case "percent:100":
        return value.toFixed(0) + "%";
      case "percent:1":
        return (value * 100).toFixed(0) + "%";
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
        return value % 1 === 0 ? value.toString() : value.toFixed(2);
    }
  }

  resize() {
    this.chart.resize();
  }
}

/**
 * Gauge chart
 */
export class GaugeChart extends Chart {
  constructor(dom: HTMLElement, info: ChartInfo) {
    super(dom, info);
  }

  protected config(): echarts.EChartsOption {
    return {
      xAxis: [{ show: false }],
    };
  }

  setData(d: any): void {
    this.chart.setOption({
      series: [
        {
          // name: d.name,
          type: 'gauge',
          // radius: '100%',
          // center: ["50%", "50%"], // 仪表位置
          // max: d.value,
          axisLine: {
            show: false,
            // lineStyle: { width: 0, opacity: 0, shadowBlur: 0 },
          },
          axisLabel: { show: false },
          axisTick: { show: false },
          splitLine: { show: false },
          pointer: { show: false },
          detail: {
            formatter: this.formatValue.bind(this),
            offsetCenter: [0, 0],
            fontSize: 64,
            fontWeight: 'bold',
          },
          data: [{ value: d.value }]
        }
      ]
    });
  }
}

/**
 * Pie chart.
 */
export class VectorChart extends Chart {
  constructor(dom: HTMLElement, info: ChartInfo) {
    super(dom, info);
  }

  protected config(): echarts.EChartsOption {
    return {
      legend: {
        show: false,
        type: 'scroll',
        align: 'right',
        orient: 'vertical',
      },
      tooltip: {
        trigger: 'item',
        formatter: (params: any): string => {
          return params.name + ": " + this.formatValue(params.value);
        },
      },
      xAxis: [{ show: false }],
      yAxis: [{ show: false }],
      series: [{
        type: this.info.type,
        radius: '90%',
        // center: ['40%', '50%'],
      }],
    };
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
  constructor(dom: HTMLElement, info: ChartInfo) {
    super(dom, info);
  }

  config(): echarts.EChartsOption {
    return {
      grid: {
        left: this.info.margin?.left || 45,
        right: this.info.margin?.right || 0,
        top: this.info.margin?.top || 10,
        bottom: this.info.margin?.bottom || 20,
      },
      legend: { show: false },
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
    };
  }

  setData(d: any): void {
    if (!d.series) {
      return;
    }

    d.series.forEach((s: any) => {
      s.type = this.info.type;
      s.areaStyle = { opacity: 0.3 };
      s.smooth = true;
      s.showSymbol = false;
      s.lineStyle = { width: 1 };
    });
    this.chart.setOption({
      legend: {
        data: d.legend,
      },
      series: d.series,
    });
  }
}

export function createChart(dom: HTMLElement, info: ChartInfo): Chart {
  switch (info.type) {
    case "gauge":
      return new GaugeChart(dom, info);
    case "line":
    case "bar":
      return new MatrixChart(dom, info);
    case "pie":
      return new VectorChart(dom, info);
    default:
      throw new Error('unknown chart type: ' + info.type)
  }
}

// export class ChartDashboardOptions {
//   name: string;
//   key?: string;
//   period?: number = 30;
//   refreshInterval?: number = 15;   // seconds
// }

// export class ChartDashboard {
//   private $panel: JQuery;
//   private readonly opts: ChartDashboardOptions;
//   private charts: Chart[] = [];
//   private timer: number;
//   private dlg: ChartDialog;

//   constructor(elem: string | Element | JQuery, charts: ChartData[], opts?: ChartDashboardOptions) {
//     this.opts = $.extend(new ChartDashboardOptions(), opts);
//     this.$panel = $(elem);
//     this.dlg = new ChartDialog(this);

//     charts.forEach(opts => this.createGraph(opts));

//     // events
//     Dispatcher.bind(this.$panel).on("remove-chart", e => {
//       let name = $(e.target).closest("div.column").data("name");
//       Modal.confirm(`Are you sure to delete chart: <strong>${name}</strong>?`, "Delete chart", dlg => {
//         this.removeGraph(name);
//         dlg.close();
//       });
//     });
//     $(window).resize(e => {
//       $.each(this.charts, (i: number, g: Chart) => {
//         g.resize();
//       });
//     });

//     this.refresh();
//   }

//   refresh() {
//     if (!this.timer) {
//       this.loadData();
//       if (this.opts.refreshInterval > 0) {
//         this.timer = setTimeout(this.refreshData.bind(this), this.opts.refreshInterval * 1000);
//       }
//     }
//   }

//   private refreshData() {
//     this.loadData();
//     if (this.opts.refreshInterval > 0) {
//       this.timer = setTimeout(this.refreshData.bind(this), this.opts.refreshInterval * 1000);
//     }
//   }

//   stop() {
//     clearTimeout(this.timer);
//     this.timer = 0;
//   }

//   setPeriod(period: number) {
//     this.opts.period = period;
//     this.loadData();
//   }

//   addGraph(opts: ChartData) {
//     this.createGraph(opts);
//     this.loadData();
//   }

//   private createGraph(opts: ChartData) {
//     for (let i = 0; i < this.charts.length; i++) {
//       let chart = this.charts[i];
//       if (chart.getOptions().name === opts.name) {
//         // chart already added.
//         return;
//       }
//     }

//     let chart = ChartFactory.create(opts);
//     if (chart != null) {
//       this.$panel.append(chart.getElem());
//       chart.init();
//       this.charts.push(chart);
//     }
//   }

//   removeGraph(name: string) {
//     let index = -1;
//     for (let i = 0; i < this.charts.length; i++) {
//       let c = this.charts[i];
//       if (c.getOptions().name === name) {
//         index = i;
//         break;
//       }
//     }

//     if (index >= 0) {
//       let $elem = this.charts[index].getElem();
//       this.charts.splice(index, 1);
//       $elem.remove();
//     }
//   }

//   save(asDefault: boolean = false) {
//     let charts: any = [];
//     this.$panel.children().each((index: number, elem: Element) => {
//       let name = $(elem).data("name");
//       for (let i = 0; i < this.charts.length; i++) {
//         let c = this.charts[i];
//         if (c.getOptions().name === name) {
//           charts.push({
//             name: c.getOptions().name,
//             width: c.getOptions().width,
//             height: c.getOptions().height,
//           });
//           break;
//         }
//       }
//     });
//     let args = {
//       name: this.opts.name,
//       key: asDefault ? '' : (this.opts.key || ''),
//       charts: charts,
//     };
//     $ajax.post(`/system/chart/save_dashboard`, args).json<AjaxResult>((r: AjaxResult) => {
//       if (r.success) {
//         Notification.show("success", "Successfully saved.");
//       } else {
//         Notification.show("danger", r.message);
//       }
//     });
//   }

//   getOptions(): ChartDashboardOptions {
//     return this.opts;
//   }

//   private loadData() {
//     if (this.charts.length == 0) {
//       return
//     }

//     let args: any = {
//       charts: this.charts.map(c => c.getOptions().name).join(","),
//       period: this.opts.period,
//     };
//     if (this.opts.key) {
//       args.key = this.opts.key;
//     }
//     $ajax.get(`/system/chart/data`, args).json((d: { [index: string]: any[] }) => {
//       $.each(this.charts, (i: number, g: Chart) => {
//         if (d[g.getOptions().name]) {
//           g.setData(d[g.getOptions().name]);
//         }
//       });
//     });
//   }
// }

// class ChartDialog {
//   private dashboard: ChartDashboard;
//   private fb: FilterBox;
//   private charts: any;
//   private $charts: JQuery;

//   constructor(dashboard: ChartDashboard) {
//     this.dashboard = dashboard;
//     this.fb = new FilterBox("#txt-query", this.filterCharts.bind(this));
//     $("#btn-add").click(this.showAddDlg.bind(this));
//     $("#btn-add-chart").click(this.addChart.bind(this));
//     $("#btn-save").click(() => {
//       this.dashboard.save();
//     });
//     $("#btn-save-as-default").click(() => {
//       this.dashboard.save(true);
//     });
//   }

//   private showAddDlg() {
//     let $panel = $("#nav-charts");
//     $panel.find("label.panel-block").remove();

//     // load charts
//     $ajax.get(`/system/chart/query`, { dashboard: this.dashboard.getOptions().name }).json((charts: any) => {
//       for (let i = 0; i < charts.length; i++) {
//         let c = charts[i];
//         $panel.append(`<label class="panel-block">
//           <input type="checkbox" value="${c.name}" data-index="${i}">${c.name}: ${c.title}
//         </label>`);
//       }
//       this.charts = charts;
//       this.$charts = $panel.find("label.panel-block");
//     });

//     let dlg = new Modal("#dlg-add-chart");
//     dlg.show();
//   }

//   private filterCharts(text: string) {
//     if (!text) {
//       this.$charts.show();
//       return;
//     }

//     this.$charts.each((i, elem) => {
//       let $elem = $(elem);
//       let texts: string[] = [
//         this.charts[i].name.toLowerCase(),
//         this.charts[i].title.toLowerCase(),
//         this.charts[i].desc.toLowerCase(),
//       ];
//       for (let i = 0; i < texts.length; i++) {
//         let index = texts[i].indexOf(text);
//         if (index >= 0) {
//           $elem.show();
//           return;
//         }
//       }
//       $elem.hide();
//     })
//   }

//   private addChart() {
//     this.$charts.each((i, e) => {
//       if ($(e).find(":checked").length > 0) {
//         let c = this.charts[i];
//         this.dashboard.addGraph(c);
//       }
//     });
//     Modal.close();
//   }
// }