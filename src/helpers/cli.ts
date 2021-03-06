import { Chapter } from "plugin";
import { isDownloaded } from "@core/downloader";
import { shutdown } from "@core/main";

import readline from 'readline';
import Table from 'cli-table3';
import cliSelect from 'cli-select';
import { TABLE_COL_WIDTHS } from "@core/constants";

const readLineAsync = (): Promise<string> => {
    const rl = readline.createInterface({
        input: process.stdin
    });

    return new Promise((resolve) => {
        rl.prompt();
        rl.on('line', (line: string) => {
            rl.close();
            resolve(line);
        });
    });
};

const getUserSelection = async <O> (values: O[]): Promise<O> => {

    let selection: any;
    try {
        selection = await cliSelect({ values: values });
    } catch (e) {
        // By definition, this is only thrown when the user sends 'SIGINT'.
        await shutdown();
    }

    return selection!.value;
};

const getUserConfirmation = async (promptString: string): Promise<String> => {
    let answerString: string;
    while (true) {
        process.stdout.write(promptString);
        answerString = await readLineAsync();

        if (answerString.toLowerCase() != "y" && answerString.toLowerCase() != "n") {
            console.log("Your answer has to be one of: (y, n).")
            continue;
        }

        break
    }

    return answerString;
}

async function getNumber(promptString: string, optChecks?: (num: number) => boolean) {
    let isNumber = (n: string) => {
        return !isNaN(parseInt(n)) && isFinite(parseInt(n));
    }

    while (true) {
        process.stdout.write(promptString);

        let numberString = await readLineAsync();

        if (isNumber(numberString) &&
        (optChecks ? optChecks(parseFloat(numberString)) : true)) {
            return parseFloat(numberString);
        }

        console.log("This number is not valid.")
    }
}

async function generateTable(chapters: Chapter[], mangaTitle: string) {
    let table = new Table({
        head: ['index', 'num', 'title', 'downloaded'],
        chars: {'mid': '', 'left-mid': '', 'mid-mid': '', 'right-mid': ''},
        colWidths: TABLE_COL_WIDTHS, wordWrap: true
    });

    chapters.forEach((item, i) => {
        const downloaded =
            isDownloaded({mangaTitle: mangaTitle, chapterTitle: item.title, num: item.num, urls: []}) ? "Y" : "N";
        table.push([chapters.length - i, item.num, item.title, downloaded]);
    })
    table.reverse();

    return table;
}

export { readLineAsync, getUserConfirmation, getNumber, generateTable, getUserSelection };