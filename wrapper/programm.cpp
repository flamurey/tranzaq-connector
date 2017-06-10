#include <windows.h>
#include <iostream>
#include <fstream>
#include <conio.h>
#include <stdio.h>

#include "XMLTranzaqConnector.h"

using namespace std;

static const int INIT_SERVER_CODE = 4;
static const int SEND_COMMAND_CODE = 1;
static const int NEW_LOG_LEVEL_CODE = 2;
static const int UNINIT_SERVER_CODE = 3;
static const int EXIT_SERVER_CODE = 9;

const BYTE LabelA[] = "!@#$%";

int readCodeMesaage() 
{
    char code = cin.get();
    return atoi(&code);
}

const int lineSize = 100000;
BYTE line[lineSize];

int readLogLevel() 
{
        char logLevel = cin.get();
        return atoi(&logLevel);
}

void sendAnswer(BYTE* answer);
void readEndl();
BYTE* readLine();

int main(int argc, char* argv[]) 
{
	XMLTranzaqConnector::getInstance();
    
    cout << "server started" << endl;
	while(true)
	{
		int logLevel;
		BYTE* answer = 0;
		int messageCode = readCodeMesaage();
		BYTE* command = 0;
                XMLTranzaqConnector* instance;
		switch (messageCode)
		{
		case INIT_SERVER_CODE:
			logLevel = readLogLevel(); 
                        readEndl();
			answer = XMLTranzaqConnector::getInstance()->init(*argv, logLevel);
			break;
		case SEND_COMMAND_CODE:
			command = readLine();
			answer = XMLTranzaqConnector::getInstance()->sendCommandToServer(command);
			break;
		case NEW_LOG_LEVEL_CODE:
			logLevel = readLogLevel();
			answer = XMLTranzaqConnector::getInstance()->setNewLogLevel(logLevel);
			readEndl();
			break;
		case UNINIT_SERVER_CODE:
			answer = XMLTranzaqConnector::getInstance()->unInitialize();
			readEndl();
			break;
		case EXIT_SERVER_CODE:
			sendAnswer(answer);
            instance =  XMLTranzaqConnector::getInstance();
			delete instance;
			return 0;
		default:
            break;
		}
		sendAnswer(answer);
	}
	return 0;
}

BYTE* readLine()
{
	int i = 0;
	BYTE ch;
	while ((ch = cin.get()) != '\n') {
		line[i] = ch;
		i++;
	}
	line[i] = 0;
	return line;
}

void readEndl() 
{
	char ch;
	do {
		ch = cin.get();
	}
	while (ch != '\n');
}

void sendAnswer(BYTE* answer)
{
	printData(LabelA);
	if (answer == 0) {
		cout << '0';
	} else {
		cout << '1';
		printData(answer);
		FreeMemory(answer);
	}
	cout << endl;
}
