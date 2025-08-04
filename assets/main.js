function loadData(url) {
    $('#loading').show();
    $.ajax({
        url: url,
        type: 'GET',
        success: function(data) {
            $('#content').html(data);
            $('#loading').hide();
        },
        error: function() {
            $('#loading').hide();
            alert('Error loading data.');
        }
    });
}

function gotoChart() {
    var sel = document.getElementById('nextChart')
    var value = sel.options[sel.selectedIndex].value;
    if(value == "") {
        return;
    }
    loadData(value);
}

function refreshChart() {
    var sd = document.getElementById('start').value;
    var ed = document.getElementById('end').value;
    loadData('/hatchets/{{.Hatchet}}/charts{{.Chart.URL}}&duration=' + sd + ',' + ed);
}

function redirect() {
    var sel = document.getElementById('table')
    var value = sel.options[sel.selectedIndex].value;
    if(value == "") {
        return;
    }
    loadData('/hatchets/' + value + '/stats/audit');
}

var input = document.getElementById("context");
if (input) {
    input.addEventListener("keypress", function(event) {
        if (event.key === "Enter") {
            event.preventDefault();
            document.getElementById("find").click();
        }
    });
}

function findLogs() {
    var sel = document.getElementById('component')
    var component = sel.options[sel.selectedIndex].value;
    sel = document.getElementById('severity')
    var severity = sel.options[sel.selectedIndex].value;
    var context = document.getElementById('context').value
    loadData('/hatchets/{{.Hatchet}}/logs/all?component='+component+'&severity='+severity+'&context='+context);
}

function getSlowopsStats() {
    var b = document.getElementById('collscan').checked;
    loadData('/hatchets/{{.Hatchet}}/stats/slowops?orderBy={{.OrderBy}}&COLLSCAN='+b);
}
function downloadStats() {
    anchor = document.createElement('a');
    anchor.download = '{{.Hatchet}}_stats.html';
    anchor.href = '/hatchets/{{.Hatchet}}/stats/slowops?type=stats&download=true';
    anchor.dataset.downloadurl = ['text/html', anchor.download, anchor.href].join(':');
    anchor.click();
}
