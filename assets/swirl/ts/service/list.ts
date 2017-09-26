///<reference path="../core/core.ts" />
namespace Swirl.Service {
    import Modal = Swirl.Core.Modal;
    import AjaxResult = Swirl.Core.AjaxResult;
    import Table = Swirl.Core.ListTable;

    export class ListPage {
        private table: Table;

        constructor() {
            this.table = new Table("#table-items");

            // bind events
            this.table.on("delete-service", this.deleteService.bind(this)).on("scale-service", this.scaleService.bind(this))
            $("#btn-delete").click(this.deleteServices.bind(this));
        }

        private deleteService(e: JQueryEventObject) {
            let $tr = $(e.target).closest("tr");
            let name = $tr.find("td:eq(1)").text().trim();
            Modal.confirm(`Are you sure to remove service: <strong>${name}</strong>?`, "Delete service", (dlg, e) => {
                $ajax.post("delete", { names: name }).trigger(e.target).encoder("form").json<AjaxResult>(r => {
                    $tr.remove();
                    dlg.close();
                })
            });
        }

        private deleteServices(e: JQueryEventObject) {
            let names = this.table.selectedKeys();
            if (names.length == 0) {
                Modal.alert("Please select one or more items.")
                return;
            }

            Modal.confirm(`Are you sure to remove ${names.length} services?`, "Delete services", (dlg, e) => {
                $ajax.post("delete", { names: names.join(",") }).trigger(e.target).encoder("form").json<AjaxResult>(r => {
                    this.table.selectedRows().remove();
                    dlg.close();
                })
            });
        }

        private scaleService(e: JQueryEventObject) {
            let $btn = $(e.target);
            let $tr = $btn.closest("tr");
            let data = {
                name: $tr.find("td:eq(1)").text().trim(),
                count: $btn.data("replicas"),
            };
            Modal.confirm(`<input name="count" value="${data.count}" class="input" placeholder="Replicas">`, "Scale service", (dlg, e) => {
                data.count = dlg.find("input[name=count]").val();
                $ajax.post("scale", data).trigger($btn).encoder("form").json<AjaxResult>(r => {
                    location.reload();
                })
            });
        }
    }
}