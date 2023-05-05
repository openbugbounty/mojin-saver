function (doc) {
    if (doc.type == "website")
        emit(doc.program, doc._id)
}