#if 1
//masala-4
#include<stdio.h>
#include<stdlib.h>
#include<time.h>

void fill_array(int n ,int m,int arr[n][m]){
    for(int i = 0;i < n;i++){
        for(int j = 0;j < m;j++){
            arr[i][j] = rand() % 26;
            printf("%d ",arr[i][j]);
        }
        puts("");
    }
    puts("");
}


void output(int n ,int m,int arr[n][m]){
    for(int i = 0;i < n;i++){
        for(int j = 0;j < m;j++){
            if(i == j || i+j == m-1){
                printf("0 ");
            }else{
                printf("%d ",arr[i][j]);
            }
        }
        puts("");
    }
    puts("");
}

int main(){
    srand(time(0));

    int n,m;
    printf(" N ta son kiriting: ");scanf("%d",&n);
    printf(" M ta son kiriting: ");scanf("%d",&m);//son kiritilganda n * m bo'ladi 

    int arr[n][m];

    fill_array(n,m,arr);
    output(n,m,arr);

    return 0;
}
#endif