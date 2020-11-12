// Arguments in order are:
//   version, bytes, name, date, email, phone
//
// A sample invocation might be:
//   deno run --allow-read --allow-write --allow-net ci/escrow/modify-pdf.ts 1.6.0-beta7 123456789 'Janice Morales' 'November 10, 2020' 'janice.morales@solo.io' '(617)-893-7557'
//
// The right way to do args parsing would be to use https://deno.land/std@0.77.0/flags
// but I could not get this to work


import { PDFDocument, rgb, degrees } from 'https://cdn.skypack.dev/pdf-lib@1.11.1?dts';

// Fetch an existing PDF file
const existingPdfBytes = await Deno.readAll(await Deno.open('ci/escrow/DepositMaterialsExhibitBFormTemplate.pdf'));

// Load an existing PDFDocument
const pdfDoc = await PDFDocument.load(existingPdfBytes);

const version = Deno.args[0]

// Put the version number in 'Deposit Name' column
const page = pdfDoc.getPage(0);
page.drawText(version, {
    x: 215,
    y: 685,
    size: 10,
    color: rgb(0.1, 0.1, 0.1),
    rotate: degrees(0),
});

// Put the number of bytes of the upload
page.drawText(Deno.args[1], {
    x: 215,
    y: 482,
    size: 10,
    color: rgb(0.1, 0.1, 0.1),
    rotate: degrees(0),
});

// Put the depositor name under 'Print Name' column
page.drawText(Deno.args[2], {
    x: 138,
    y: 268,
    size: 10,
    color: rgb(0.1, 0.1, 0.1),
    rotate: degrees(0),
});

// Put the deposit date
page.drawText(Deno.args[3], {
    x: 138,
    y: 255,
    size: 10,
    color: rgb(0.1, 0.1, 0.1),
    rotate: degrees(0),
});

// Put the depositor contact email under 'Email Address'
page.drawText(Deno.args[4], {
    x: 138,
    y: 243,
    size: 10,
    color: rgb(0.1, 0.1, 0.1),
    rotate: degrees(0),
});

// Put the depositor contact phone number under 'Telephone Number'
page.drawText(Deno.args[5], {
    x: 138,
    y: 230,
    size: 10,
    color: rgb(0.1, 0.1, 0.1),
    rotate: degrees(0),
});

// Save the PDFDocument and write it to a file
const pdfBytes = await pdfDoc.save();
const pdfFileName = '_gloo-ee-source/ExhibitB' + version + '.pdf'
await Deno.writeFile(pdfFileName, pdfBytes);

// Done! 💥
console.log('PDF file written to ' + pdfFileName);