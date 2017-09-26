///<reference path="../core/core.ts" />
namespace Swirl.Service {
    import Validator = Swirl.Core.Validator;
    import OptionTable = Swirl.Core.OptionTable;
    import EditTable = Swirl.Core.EditTable;
    import Modal = Swirl.Core.Modal;
    import Table = Swirl.Core.Table;

    class ServiceModeRule implements Swirl.Core.ValidationRule {
        private $mode: JQuery;

        constructor($model: JQuery) {
            this.$mode = $model;
        }

        validate($form: JQuery, $input: JQuery, arg?: string): { ok: boolean, error?: string } {
            if (this.$mode.val() == "global") {
                return { ok: true }
            }

            const regex = /^[1-9]\d*$/;
            return { ok: regex.test($.trim($input.val())) };
        }
    }

    class MountTable extends EditTable {
        protected render(): string {
            return `<tr>
              <td>
                <div class="select is-small">
                  <select name="mounts[${this.index}].type">
                    <option value="bind">Bind</option>
                    <option value="volume">Volume</option>
                    <option value="tmpfs">TempFS</option>
                  </select>
                </div>
              </td>
              <td>
                <input name="mounts[${this.index}].src" class="input is-small" placeholder="path in host">
              </td>
              <td><input name="mounts[${this.index}].dst" class="input is-small" placeholder="path in container"></td>
              <td>
                <div class="select is-small">
                  <select name="mounts[${this.index}].read_only" data-type="bool">
                    <option value="false">No</option>
                    <option value="true">Yes</option>
                  </select>
                </div>
              </td>
              <td>
                <div class="select is-small">
                  <select name="mounts[${this.index}].propagation">
                    <option value="">--Select--</option>
                    <option value="rprivate">rprivate</option>
                    <option value="private">private</option>
                    <option value="rshared">rshared</option>
                    <option value="shared">shared</option>
                    <option value="rslave">rslave</option>
                    <option value="slave">slave</option>
                  </select>
                </div>
              </td>
              <td>
                <a class="button is-small is-outlined is-danger" data-action="delete-mount">
                  <span class="icon is-small">
                    <i class="fa fa-trash"></i>
                  </span>
                </a>
              </td>
            </tr>`;
        }
    }

    class PortTable extends EditTable {
        protected render(): string {
            return `<tr>
                <td><input name="endpoint.ports[${this.index}].published_port" class="input is-small" placeholder="port in host" data-type="integer"></td>
                <td>
                  <input name="endpoint.ports[${this.index}].target_port" class="input is-small" placeholder="port in container" data-type="integer">
                </td>
                <td>
                  <div class="select is-small">
                    <select name="endpoint.ports[${this.index}].protocol">
                      <option value="false">TCP</option>
                      <option value="true">UDP</option>
                    </select>
                  </div>
                </td>
                <td>
                  <div class="select is-small">
                    <select name="endpoint.ports[${this.index}].publish_mode">
                      <option value="rprivate">ingress</option>
                      <option value="private">host</option>
                    </select>
                  </div>
                </td>
                <td>
                  <a class="button is-small is-outlined is-danger" data-action="delete-endpoint-port">
                    <span class="icon is-small">
                      <i class="fa fa-trash"></i>
                    </span>
                  </a>
                </td>
              </tr>`;
        }
    }

    class ConstraintTable extends EditTable {
        protected render(): string {
            return `<tr>
                <td>
                  <input name="placement.constraints[${this.index}].name" class="input is-small" placeholder="e.g. node.role/node.hostname/node.id/node.labels.*/engine.labels.*/...">
                </td>
                <td>
                  <div class="select is-small">
                    <select name="placement.constraints[${this.index}].op">
                      <option value="==">==</option>
                      <option value="!=">!=</option>
                    </select>
                  </div>
                </td>
                <td>
                  <input name="placement.constraints[${this.index}].value" class="input is-small" placeholder="e.g. manager">
                </td>
                <td>
                  <a class="button is-small is-outlined is-danger" data-action="delete-constraint">
                    <span class="icon is-small">
                      <i class="fa fa-trash"></i>
                    </span>
                  </a>
                </td>
              </tr>`;
        }
    }

    class PreferenceTable extends EditTable {
        protected render(): string {
            return `<tr>
                <td>
                  <input name="placement.preferences[${this.index}].spread" class="input is-small" placeholder="e.g. engine.labels.az">
                </td>
                <td>
                  <a class="button is-small is-outlined is-danger" data-action="delete-preference">
                    <span class="icon is-small">
                      <i class="fa fa-trash"></i>
                    </span>
                  </a>
                </td>
              </tr>`;
        }
    }

    class ConfigTable extends Table {
        public readonly name: string;
        private index: number;
        private $body: JQuery;

        constructor(elem: string | JQuery | Element) {
            super(elem);

            this.name = this.$table.data("name");
            this.$body = this.$table.find("tbody");
            this.index = this.$body.find("tr").length;

            super.on("add-" + this.name, this.showAddDialog.bind(this)).on("delete-" + this.name, ConfigTable.deleteRow);
        }

        public addRow(id: string, name: string) {
            let field = `${this.name}s[${this.index}]`;
            this.$body.append(`<tr>
                <td>${name}<input name="${field}.id" value="${id}" type="hidden"><input name="${field}.name" value="${name}" type="hidden"></td>
                <td><input name="${field}.file_name" value="${name}" class="input is-small"></td>
                <td><input name="${field}.uid" value="0" class="input is-small"></td>
                <td><input name="${field}.gid" value="0" class="input is-small"></td>
                <td><input name="${field}.mode" value="444" class="input is-small" data-type="integer"></td>
                <td>
                  <a class="button is-small is-outlined is-danger" data-action="delete-${this.name}">
                    <span class="icon is-small">
                      <i class="fa fa-remove"></i>
                    </span>
                  </a>
                </td>
              </tr>`);
            this.index++;
        }

        private showAddDialog(e: JQueryEventObject) {
            let dlg = new Modal("#dlg-add-"+this.name);
            dlg.find(":checked").prop("checked", false);
            dlg.error();
            dlg.show();
        }

        private static deleteRow(e: JQueryEventObject) {
            $(e.target).closest("tr").remove();
        }
    }

    export class EditPage {
        private $mode: JQuery;
        private $replicas: JQuery;
        private secret: ConfigTable;
        private config: ConfigTable;

        constructor() {
            this.$mode = $("#cb-mode");
            this.$replicas = $("#txt-replicas");
            new OptionTable("#table-envs");
            new OptionTable("#table-slabels");
            new OptionTable("#table-clabels");
            new OptionTable("#table-log_driver-options");
            new PortTable("#table-endpoint-ports");
            new MountTable("#table-mounts");
            new ConstraintTable("#table-constraints");
            new PreferenceTable("#table-preferences");
            this.secret = new ConfigTable("#table-secrets");
            this.config = new ConfigTable("#table-configs");

            // register custom validators
            Validator.register("service-mode", new ServiceModeRule(this.$mode), "Please input a positive integer.");

            // bind events
            this.$mode.change(e => this.$replicas.toggle(this.$mode.val() != "global"))
            $("#btn-add-secret").click(() => EditPage.addConfig(this.secret));
            $("#btn-add-config").click(() => EditPage.addConfig(this.config));
        }

        private static addConfig(t: ConfigTable) {
            let dlg = Modal.current();
            let $cbs = dlg.find(":checked");
            if ($cbs.length == 0) {
                dlg.error(`Please select the ${t.name} files.`)
            } else {
                dlg.close();
                $cbs.each((i, cb) => {
                    let $cb = $(cb);
                    t.addRow($cb.val(), $cb.data("name"));
                });
            }
        }
    }

    export class NewPage extends EditPage {
        private $registry: JQuery;
        private $registryUrl: JQuery;

        constructor() {
            super();

            this.$registryUrl = $("#a-registry-url");
            this.$registry = $("#cb-registry");
            this.$registry.change(e => {
                let url = this.$registry.find("option:selected").data("url") || "";
                this.$registryUrl.text(url + "/").toggle(url != "");
            })
        }
    }
}