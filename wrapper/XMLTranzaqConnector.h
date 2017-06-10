#include <windows.h>
#include <iostream>
#include <fstream>
#include <algorithm>

using namespace std;

#define TACCESS_API  __declspec(dllimport)
typedef bool (WINAPI *tcallback)(BYTE* pData);
typedef BYTE* (WINAPI *typeSendCommand)(BYTE* pData);
typedef bool (WINAPI  *typeFreeMemory)(BYTE* pData);
typedef bool (WINAPI *typeSetCallback)(tcallback pCallback);
typedef BYTE* (WINAPI *typeInitialize)(const BYTE* dir, int level);
typedef BYTE* (WINAPI *typeSetLogLevel)(int level);
typedef BYTE* (WINAPI *typeUninitialize)();



typeFreeMemory FreeMemory;
const BYTE LabelB[] = "%$#@!";


void printData(const BYTE* data) 
{
	BYTE ch;
	int i = 0;
	while ( (ch= *(data+i))!= 0)
	{
		i++;
                if (ch != '\n' && ch != '\r')
			std::cout << ch;
	}
}

bool CALLBACK acceptor(BYTE *pData)
{
	printData(LabelB);
	printData(pData);
	std::cout<<std::endl;
	FreeMemory(pData);
	return true;
}

class XMLTranzaqConnector
{
private:
	static XMLTranzaqConnector* instance;

	typeInitialize Initialize;
	typeUninitialize UnInitialize;
	typeSetLogLevel SetLogLevel;
	typeSetCallback SetCallback;
	typeSendCommand SendCommand;

	HMODULE hm;
	XMLTranzaqConnector()
	{
		loadLibrary();
		SetCallback(acceptor);
	}

	

	void loadLibrary()
	{
		setlocale(LC_CTYPE, "");
		hm = LoadLibrary("txmlconnector64.dll");
		if (hm) {
            Initialize=  reinterpret_cast<typeInitialize>(GetProcAddress(hm, "Initialize"));

            UnInitialize = reinterpret_cast<typeUninitialize>(GetProcAddress(hm, "UnInitialize"));

            SetLogLevel = reinterpret_cast<typeSetLogLevel>(GetProcAddress(hm, "SetLogLevel"));

            SetCallback = reinterpret_cast<typeSetCallback>(GetProcAddress(hm, "SetCallback"));

            SendCommand = reinterpret_cast<typeSendCommand>(GetProcAddress(hm,"SendCommand"));

            FreeMemory =
                reinterpret_cast<typeFreeMemory>(GetProcAddress(hm, "FreeMemory"));
		} else {
            std::cout << "Fail in LoadLibrary"<< std::endl;
        }
	}

	

public:

	~XMLTranzaqConnector()
	{
        unloadLibrary();
	}

	void unloadLibrary()
	{
		try {
			FreeLibrary(hm);
		}
		catch (...) {
			std::cout<<"Fail in FreeLibrary"<<std::endl;
		}
	}
	

	static XMLTranzaqConnector* getInstance() {
		if (instance == 0)
			instance = new XMLTranzaqConnector();
		return instance;
	}

	BYTE* init(char* exePath,int logLevel) 
	{
                char addingPath[] = "tc_logs";
		char logPath[200];
		char ch;
		int len = 0;
		int lastSlah = -1;
		while ( (ch = *(exePath + len)) != 0)
		{
			logPath[len] = ch;
			if (ch == '\\')
				lastSlah = len;
			len++;
		}
		int continueIndex;
		if (lastSlah != -1)
			continueIndex = lastSlah + 1;
		else
			continueIndex = 0;
		len = 0;
		while ( (ch = *(addingPath + len)) != 0)
		{
			logPath[continueIndex] = ch;
			continueIndex++;
			len++;
		}
                logPath[continueIndex] = 0;

                //const WCHAR* wLogPath = convertCharToWchar(logPath);
                CreateDirectory(logPath, NULL);
                //delete [] wLogPath;

		const BYTE* path = reinterpret_cast<const BYTE*>(logPath);
		BYTE* res = Initialize(path, logLevel);
		return res;

		
	}

        const WCHAR* convertCharToWchar(char* input) {
            const WCHAR *res;
            // required size
            int nChars = MultiByteToWideChar(CP_ACP, 0, input, -1, NULL, 0);
            // allocate it
            res = new WCHAR[nChars];
            MultiByteToWideChar(CP_ACP, 0, input, -1, (LPWSTR)res, nChars);
            return res;
        }


	BYTE* sendCommandToServer(BYTE* command) 
	{
		return SendCommand(command);
	}
	
	BYTE* setNewLogLevel(int newLevel)
	{
		return SetLogLevel(newLevel);
	}

	BYTE* unInitialize()
	{
		return UnInitialize();
	}
};


XMLTranzaqConnector* XMLTranzaqConnector::instance = 0;
