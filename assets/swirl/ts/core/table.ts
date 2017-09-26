/*!
 * Swirl Table Library v1.0.0
 * Copyright 2017 cuigh. All rights reserved.
 *
 * @author cuigh(noname@live.com)
 */
///<reference path="dispatcher.ts"/>
namespace Swirl.Core {
    export class Table {
        protected $table: JQuery;
        private dispatcher: Dispatcher;

        constructor(table: string | Element | JQuery) {
            this.$table = $(table);
            this.dispatcher = Dispatcher.bind(this.$table);
        }

        /**
         * Bind action
         *
         * @param action
         * @param handler
         * @returns {Swirl.Core.ListTable}
         */
        on(action: string, handler: (e: JQueryEventObject) => any): this {
            this.dispatcher.on(action, handler);
            return this;
        }
    }

    export class ListTable extends Table {
        constructor(table: string | Element | JQuery) {
            super(table);

            this.on("check-all", e => {
                let checked = (<HTMLInputElement>e.target).checked;
                this.$table.find("tbody>tr").each((i, elem) => {
                    $(elem).find("td:first>:checkbox").prop("checked", checked);
                });
            });
            this.on("check", () => {
                let rows = this.$table.find("tbody>tr").length;
                let checkedRows = this.selectedRows().length;
                this.$table.find("thead>tr>th:first>:checkbox").prop("checked", checkedRows > 0 && rows == checkedRows);
            })
        }

        /**
         * Return selected rows.
         */
        selectedRows(): JQuery {
            return this.$table.find("tbody>tr").filter((i, elem) => {
                let cb = $(elem).find("td:first>:checkbox");
                return cb.prop("checked");
            });
        }

        /**
         * Return keys of selected items.
         */
        selectedKeys(): string[] {
            let keys: string[] = [];
            this.$table.find("tbody>tr").each((i, elem) => {
                let cb = $(elem).find("td:first>:checkbox");
                if (cb.prop("checked")) {
                    keys.push(cb.val());
                }
            });
            return keys;
        }
    }

    export abstract class EditTable extends Table {
        protected name: string;
        protected index: number;
        protected alias: string;

        constructor(elem: string | JQuery | Element) {
            super(elem);

            this.name = this.$table.data("name");
            this.alias = this.name.replace(".", "-");
            this.index = this.$table.find("tbody>tr").length;

            super.on("add-" + this.alias, this.addRow.bind(this)).on("delete-" + this.alias, OptionTable.deleteRow);
        }

        protected abstract render(): string;

        private addRow() {
            this.$table.find("tbody").append(this.render());
            this.index++;
        }

        private static deleteRow(e: JQueryEventObject) {
            $(e.target).closest("tr").remove();
        }
    }

    export class OptionTable extends EditTable {
        constructor(elem: string | JQuery | Element) {
            super(elem);
        }

        protected render(): string {
            return `<tr>
            <td><input name="${this.name}s[${this.index}].name" class="input is-small" type="text"></td>
            <td><input name="${this.name}s[${this.index}].value" class="input is-small" type="text"></td>
            <td>
              <a class="button is-small is-danger is-outlined" data-action="delete-${this.alias}">
                <span class="icon is-small"><i class="fa fa-trash"></i></span>
              </a>            
            </td>
          </tr>`;
        }
    }
}