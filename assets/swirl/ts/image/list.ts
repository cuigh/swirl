///<reference path="../core/core.ts" />
namespace Swirl.Image {
    import Modal = Swirl.Core.Modal;
    import AjaxResult = Swirl.Core.AjaxResult;
    import Table = Swirl.Core.ListTable;

    export class ListPage {
        private table: Table;

        constructor() {
            this.table = new Table("#table-items");

            // bind events
            this.table.on("delete-image", this.deleteImage.bind(this));
            $("#btn-delete").click(this.deleteImages.bind(this));
        }

        private deleteImage(e: JQueryEventObject) {
            let $tr = $(e.target).closest("tr");
            let name = $tr.find("td:eq(1)").text().trim();
            let id = $tr.find(":checkbox:first").val();
            Modal.confirm(`Are you sure to remove image: <strong>${name}</strong>?`, "Delete image", (dlg, e) => {
                $ajax.post("delete", { ids: id }).trigger(e.target).encoder("form").json<AjaxResult>(() => {
                    $tr.remove();
                    dlg.close();
                })
            });
        }

        private deleteImages() {
            let ids = this.table.selectedKeys();
            if (ids.length == 0) {
                Modal.alert("Please select one or more items.");
                return;
            }

            Modal.confirm(`Are you sure to remove ${ids.length} images?`, "Delete images", (dlg, e) => {
                $ajax.post("delete", { ids: ids.join(",") }).trigger(e.target).encoder("form").json<AjaxResult>(() => {
                    this.table.selectedRows().remove();
                    dlg.close();
                })
            });
        }
    }
}