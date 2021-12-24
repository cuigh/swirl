import * as echarts from "echarts";
import { ChartInfo } from "@/api/dashboard";
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
