function (doc) {
    for (prop in doc.tags)
        emit([doc.program ? doc.program : doc._id, doc.type], prop);
}