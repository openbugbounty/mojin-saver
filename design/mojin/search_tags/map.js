function (doc) {
    if (doc.tags) {
        for (prop in doc.tags) {
            if (Array.isArray(doc.tags[prop])) {
                for (val in doc.tags[prop]) {
                    emit([prop, doc.tags[prop][val]], [doc.type, doc._id, doc.program ? doc.program : doc._id])
                }
            } else {
                emit([prop, doc.tags[prop]], [doc.type, doc._id, doc.program ? doc.program : doc._id])
            }
        }
    }

    if (doc.type == 'website') {
        emit(['port', doc.port + ''], [doc.type, doc._id, doc.program])
        emit(['hostname', doc.hostname], [doc.type, doc._id, doc.program])
        emit(['scheme', doc.scheme], [doc.type, doc._id, doc.program])
    }

    if (doc.source) {
        emit(['source', doc.source], [doc.type, doc._id, doc.program])
    }
}